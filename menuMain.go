package main

import (
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/caseymrm/menuet"
)

// menu items
func menuItems() []item {
	checkDictionaries()

	// getRandomEmoji()
	currCreds.uname.value = generateUsername()
	currCreds.pass.value = generatePassword()
	clipboard.WriteAll(currCreds.pass.value)
	wow := generateWoW()

	spacer := item{}

	return []item{
		item{
			Text:       passwordMachineText(),
			FontWeight: menuet.WeightBold,
			FontSize:   18,
			Clicked: func() {
				debugNotification("clicked Title")
			},
		},
		item{},
		item{Text: "Username"},
		menuDisplayCredential(currCreds.uname.value, "username"),
		item{Text: "Password (" + strconv.Itoa(len(currCreds.pass.value)) + " characters)"},
		menuDisplayCredential(currCreds.pass.value, "password"),
		spacer,
		submenuLastItemClicked(),
		spacer,
		item{
			Text: "Words of Wisdom"},
		menuDisplayCredential(wow, "none"),
		spacer,
		submenuSettings(),
	}
}

func menuDisplayCredential(details, mode string) item {
	return item{
		Text:       details,
		FontWeight: menuet.WeightMedium,
		FontSize:   16,
		Clicked: func() {
			clipboard.WriteAll(details)
			switch mode {
			case "username":
				clickedCreds.uname.value = details
				// clickedCreds.uname.time = time.Now()
			case "password":
				clickedCreds.pass.value = details
				// clickedCreds.pass.time = time.Now()
			}
			setDefaults()
			// credentials.update()
		},
	}
}

func passwordMachineText() string {
	if !config.devVersion {
		return "Password machine"
	}
	return "Password machine\t" + getEmojisSeeded(2, executableLastModTime())
}
