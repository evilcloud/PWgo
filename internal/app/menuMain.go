package app

import (
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/caseymrm/menuet"
)

// menu items
func (a *App) menuItems() []item {
	// getRandomEmoji()
	a.currCreds.uname.value = a.generateUsername()
	a.currCreds.pass.value = a.generatePassword()
	clipboard.WriteAll(currCreds.pass.value)
	wow := a.generateWoW()

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

func (a *App) menuDisplayCredential(details, mode string) item {
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

func (a *App) passwordMachineText() string {
	if !a.config.devVersion {
		return "Password machine"
	}
	return "Password machine\t" + a.emojisSeeded(2, executableLastModTime())
}
