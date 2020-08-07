package app

import (
	"strconv"

	"github.com/caseymrm/menuet"
)

func (a *App) submenuLastItemClicked() menuet.MenuItem {
	return item{
		Text: "Last copy-clicked",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				item{Text: "Username"},
				item{
					Text:     humaniseDuration(a.clickedCreds.uname.time),
					FontSize: 10,
				},
				menuDisplayCredential(a.clickedCreds.uname.value, "clicked"),
				item{},
				item{Text: "Password (" + strconv.Itoa(len(a.clickedCreds.Pass.value)) + " characters)"},
				item{
					Text:     humaniseDuration(a.clickedCreds.Pass.time),
					FontSize: 10,
				},
				menuDisplayCredential(a.clickedCreds.Pass.value, "clicked"),
			}
		},
	}
}
