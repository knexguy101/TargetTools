package target

import (
	"TargetTools/console"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/stealth"
	"time"
)

type GenAccount struct {
	Username string
	Password string
	Firstname string
	Lastname string
}

func Generate(genAccount *GenAccount) error {
	console.Write(fmt.Sprintf("Starting Generator [%s]", genAccount.Username))
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(fmt.Sprintf("Fatal Error: %s", err))
			return
		}
	}()

	// Create a new page
	console.Write("Navigating")
	page := stealth.MustPage(browser)
	page.MustNavigate("https://www.target.com/")
	if err = page.Navigate("https://gsp.target.com/gsp/authentications/v1/auth_codes?client_id=ecom-web-1.0.0&state=1621662171222&redirect_uri=https%3A%2F%2Fwww.target.com%2F&assurance_level=M"); err != nil {
		return err
	}

	console.Write("Creating Account")
	page.MustElement("#createAccount").MustClick()

	console.Write("Entering Details")
	page.MustElement("#username").MustInput(genAccount.Username)
	page.MustElement("#password").MustInput(genAccount.Password)
	page.MustElement("#firstname").MustInput(genAccount.Firstname)
	page.MustElement("#lastname").MustInput(genAccount.Lastname)

	time.Sleep(2 * time.Second)

	console.Write("Confirming Account")
	page.MustElement("#createAccount").MustClick()
	page.MustWaitNavigation()

	page.MustElement("#circle-skip").MustClick()
	page.MustWaitNavigation()

	console.Write("Created Account")

	return nil
}