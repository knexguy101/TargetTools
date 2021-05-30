package target

import (
	"TargetTools/console"
	"TargetTools/gmail"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/stealth"
	"time"
)

func Reset(username, password string, gmailHandler *gmail.GmailHandler) error {

	console.Write(fmt.Sprintf("Resetting %s", username))

	var err error
	browser := rod.New().DefaultDevice(devices.Device {
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
		Screen: devices.Screen{
			Horizontal: devices.ScreenSize{
				Width: 1920,
				Height: 1080,
			},
			Vertical: devices.ScreenSize{
				Width: 1920,
				Height: 1080,
			},
		},
		AcceptLanguage: "en-US",
	}).MustConnect()

	// Even you forget to close, rod will close it after main process ends.
	defer browser.MustClose()

	// Create a new page
	console.Write("Navigating")
	page := stealth.MustPage(browser)
	page.MustNavigate("https://www.target.com/")
	if err = page.Navigate("https://gsp.target.com/gsp/authentications/v1/auth_codes?client_id=ecom-web-1.0.0&state=1621662171222&redirect_uri=https%3A%2F%2Fwww.target.com%2F&assurance_level=M"); err != nil {
		return err
	}

	console.Write("Recover")
	page.MustElement("#recoveryPassword").MustClick()

	console.Write("Entering Username")
	page.MustElement("#username").MustInput(username)
	page.MustElement("#continue").MustClick()
	page.MustElement("#continue").MustClick()

	console.Write("Waiting for Email")
	var newEmails *[]string
	for true {
		newEmails, err = gmailHandler.RefreshEmails()
		if err != nil {
			fmt.Println(err)
			time.Sleep(2 * time.Second)
		} else if len(*newEmails) > 0 {
			break
		}
	}
	val := (*newEmails)[0]

	console.Write(fmt.Sprintf("Confirmation Code: %s", val))
	page.MustElement("input[type='tel']").MustInput(val).MustClick()
	time.Sleep(1 * time.Second)
	page.MustElement("#verify").MustClick()

	time.Sleep(7 * time.Second)

	console.Write("Entering New Password")
	page.MustElement("#password").MustInput(password)
	time.Sleep(1 * time.Second)
	page.MustElement("#submit").MustClick()

	time.Sleep(5 * time.Second)
	console.Write("Submittited Password")

	return nil
}