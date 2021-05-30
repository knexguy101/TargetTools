package main

import (
	"TargetTools/console"
	"TargetTools/file"
	"TargetTools/gmail"
	"TargetTools/target"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"log"
)

var gmailHandler *gmail.GmailHandler

func main() {

	gmailHandler = gmail.NewGmailHandler()
	gmailHandler.Login()
	gmailHandler.RefreshEmails() //get our initial emails

	console.WriteMultipleLines("Options", []string {
		"[1] Reset",
		"[2] Generate",
		"[3] Help",
	})
	line, err := console.ReadPromptedLine("Pick an option")
	if err != nil {
		log.Fatalln(err) //how
	}
	switch line {
	case "1":
		handleReset()
		break
	case "2":
		handleGenerate()
		break
	case "3":
		handleHelp()
		break
	}
}

func handleReset() {
	accountPath, err := console.ReadPromptedLine("Drag and drop the file containing the accounts")
	if err != nil {
		log.Fatalln(err)
	}
	accountData, err := file.ReadLines(accountPath)
	if err != nil {
		log.Fatalln(err)
	}

	newPassword, err := console.ReadPromptedLine("Enter the new password")
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range accountData {
		err = target.Reset(v, newPassword, gmailHandler)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handleGenerate() {
	accountPath, err := console.ReadPromptedLine("Drag and drop the file containing the accounts")
	if err != nil {
		log.Fatalln(err)
	}
	accountData, err := file.ReadLines(accountPath)
	if err != nil {
		log.Fatalln(err)
	}

	password, err := console.ReadPromptedLine("Enter the password")
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range accountData {
		err = target.Generate(&target.GenAccount{
			Firstname: randomdata.FirstName(randomdata.Male),
			Lastname: randomdata.LastName(),
			Username: v,
			Password: password,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handleHelp() {
	gmail.OpenBrowser("https://github.com/knexguy101/TargetTools") //im lazy
}
