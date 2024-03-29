package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

// Configs defines configuraiton of this app
type Configs struct {
	TriggerMenus    []string `json:"triggerMenus"`
	MDmealURL       string   `json:"mdmealURL"`
	MDmealAccount   User     `json:"mdmealAcount"`
	LINENotifyToken string   `json:"lineNotifyToken"`
	MaxRetryCount   int      `json:"maxRetryCount"`
}

// User defines user info at mdmeal
type User struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// Menu defines menu info
type Menu struct {
	Name string
	Date string
}

type retryableErr struct {
	err error
}

func (e *retryableErr) Error() string {
	return e.err.Error()
}

func (e *retryableErr) UnWrap() error {
	return e.err
}

// Notification defines type of a notify function
type Notification func(string) error

func readConfigs(path string) (*Configs, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var configs Configs
	if err := json.Unmarshal(content, &configs); err != nil {
		return nil, err
	}
	return &configs, nil
}

func downloadMenu(mdmealURL string, user *User) (io.Reader, error) {
	chromeArgs := agouti.ChromeOptions(
		"args", []string{
			"--headless",
			"--disable-gpu",
		})
	chromeExcludeSwitches := agouti.ChromeOptions(
		"excludeSwitches", []string{
			"enable-logging",
		})

	driver := agouti.ChromeDriver(chromeArgs, chromeExcludeSwitches)
	defer driver.Stop()
	if err := driver.Start(); err != nil {
		return nil, err
	}

	page, err := driver.NewPage()
	if err != nil {
		return nil, err
	}
	if err := page.Session().SetImplicitWait(3000); err != nil {
		return nil, err
	}
	if err := page.Session().SetPageLoad(3000); err != nil {
		return nil, err
	}
	fmt.Println("Opening the welcome page...")
	if err := page.Navigate(mdmealURL); err != nil {
		return nil, err
	}
	fmt.Println("Logging in...")
	if err := page.FindByID("txtLoginId").SendKeys(user.ID); err != nil {
		return nil, &retryableErr{err}
	}
	if err := page.FindByID("txtPassword").SendKeys(user.Password); err != nil {
		return nil, &retryableErr{err}
	}
	if err := page.FindByID("ibLogin").Click(); err != nil {
		return nil, &retryableErr{err}
	}

	fmt.Println("Opening the menu page...")
	time.Sleep(time.Millisecond * 100)
	if err := page.FindByID("ibOrder").Click(); err != nil {
		return nil, err
	}

	fmt.Println("Downloading menus...")
	time.Sleep(time.Millisecond * 100)
	page.FindByID("gvOrder")

	html, err := page.HTML()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader([]byte(html)), nil
}

func scrape(html io.Reader) ([]Menu, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, err
	}
	orderHTML := doc.Find("table#gvOrder > tbody > tr:nth-child(n+1) > td:nth-child(5) > div.meal-meau-title")
	if orderHTML.Length() == 0 {
		return nil, errors.New("scraped result is empty")
	}
	menuNames := make([]string, orderHTML.Length())
	orderHTML.Each(func(idx int, s *goquery.Selection) {
		menuNames[idx] = strings.TrimSpace(s.Text())
	})
	dateHTML := doc.Find("table#gvOrder > tbody > tr:nth-child(n+1) > td:nth-child(1)")
	menuDates := make([]string, dateHTML.Length())
	dateHTML.Each(func(idx int, s *goquery.Selection) {
		menuDates[idx] = s.Text()
	})

	menus := []Menu{}
	for i := 0; i < len(menuNames); i++ {
		if menuNames[i] != "" {
			menus = append(menus, Menu{menuNames[i], menuDates[i]})
		}
	}
	return menus, nil
}

// func loadHTMLFromFile(cacheFilePath string) (io.Reader, error) {
// 	content, err := ioutil.ReadFile(cacheFilePath)
// 	return bytes.NewReader(content), err
// }

func notifyToLINE(msg, token string) error {
	values := url.Values{}
	values.Add("message", msg)

	req, err := http.NewRequest(
		"POST",
		"https://notify-api.line.me/api/notify",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func notifyErrorAndExit(err error, notify Notification) {
	notify(err.Error())
	log.Fatalln(err)
}

func main() {
	var notify Notification
	var configFilePath string

	flag.Parse()
	if len(flag.Args()) == 1 {
		configFilePath = flag.Args()[0]
	} else {
		configFilePath = "configs.json"
	}
	configs, err := readConfigs(configFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	notify = func(msg string) error { return notifyToLINE(msg, configs.LINENotifyToken) }
	fmt.Println("Downloading html...")

	var html io.Reader
	for cnt := 0; cnt < configs.MaxRetryCount; cnt++ {
		var errRetryable *retryableErr
		html, err = downloadMenu(configs.MDmealURL, &configs.MDmealAccount)
		if err == nil {
			break
		}
		if ok := errors.As(err, &errRetryable); !ok || cnt == configs.MaxRetryCount-1 {
			notifyErrorAndExit(err, notify)
		} else {
			fmt.Println("Retry after 3 seconds")
			time.Sleep(time.Second * 3)
		}
	}

	fmt.Println("Scraping html...")
	menus, err := scrape(html)
	if err != nil {
		notifyErrorAndExit(err, notify)
	}
	notified := false
	for _, menu := range menus {
		for _, triggerWord := range configs.TriggerMenus {
			if strings.Contains(menu.Name, triggerWord) {
				msg := fmt.Sprintf("%v@%v", menu.Name, menu.Date)
				fmt.Println(msg)
				if err := notify(msg); err != nil {
					log.Fatalln(err)
				}
				notified = true
				break
			}
		}
	}
	if !notified {
		fmt.Println("Finished without notifications")
	}
}
