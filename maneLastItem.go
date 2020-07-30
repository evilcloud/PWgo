package main

import (
	"strconv"

	"github.com/caseymrm/menuet"
)

func submenuLastItemClicked() menuet.MenuItem {
	return item{
		Text: "Last copy-clicked",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				item{Text: "Username"},
				item{
					Text:     humaniseDuration(clickedCreds.uname.time),
					FontSize: 10,
				},
				menuDisplayCredential(clickedCreds.uname.value, "clicked"),
				item{},
				item{Text: "Password " + strconv.Itoa(config.passLength) + " characters)"},
				item{
					Text:     humaniseDuration(clickedCreds.pass.time),
					FontSize: 10,
				},
				menuDisplayCredential(clickedCreds.pass.value, "clicked"),
			}
		},
	}
}
