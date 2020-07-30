package main

import (
	"strconv"

	"github.com/caseymrm/menuet"
)

func submenuSettings() item {
	subSubLengthItem := func(length int) item {
		return item{
			Text: strconv.Itoa(length),
			Clicked: func() {
				config.passLength = length
			},
			State: config.passLength == length,
		}
	}
	setDefaults()
	return item{
		Text: "Settings",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				{Text: "Password length"},
				subSubLengthItem(passShort),
				subSubLengthItem(passAcceptable),
				subSubLengthItem(passStandard),
				subSubLengthItem(passLong),
				item{},
				// {Text: "Additional security"},
				// submenuAdditionalSecurity(),
				// item{},
				item{Text: "Level of profanity"},
				nsfwItem(),
				sailorItem(),
				subSubVersion(),
			}
		},
	}
}

func sailorItem() menuet.MenuItem {
	sailorItem := item{Text: "Sailor-redneck mode (only in NSFW mode)"}
	if config.profanity.nsfw {
		sailorItem = item{Text: "Sailor-redneck mode",
			Clicked: func() {
				config.loadDict = false
				if config.profanity.sailor {
					config.profanity.sailor = false
					config.profanity.nsfw = true
				} else {
					config.profanity.sailor = true
					menuet.App().Notification(menuet.Notification{
						Title:   "A less secure novelty setting.",
						Message: "Also using it will make you look like a juvenile asshole. Use at your own risk.",
					})
				}
				setMenuState()
			},
			State: config.profanity.sailor,
		}
	}
	return sailorItem
}

func nsfwItem() menuet.MenuItem {
	return item{
		Text: "NSFW",
		Clicked: func() {
			config.loadDict = false
			if config.profanity.nsfw {
				config.profanity.sfw = true
				config.profanity.nsfw = false
			} else {
				config.profanity.sfw = false
				config.profanity.nsfw = true
			}
			config.profanity.sailor = false
			setMenuState()
		},
		State: config.profanity.nsfw}
}

func subSubVersion() item {
	return item{
		Text:     Version,
		FontSize: 7,
	}
}
