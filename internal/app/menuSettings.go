package app

import (
	"strconv"

	"github.com/caseymrm/menuet"
)

func (a *App) submenuSettings() item {
	subSubLengthItem := func(length int) item {
		return item{
			Text: strconv.Itoa(length),
			Clicked: func() {
				a.config.PassLength = length
			},
			State: a.config.PassLength == length,
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
				{Text: "Additional security"},
				submenuAdditionalSecurity(),
				item{},
				item{Text: "Level of profanity"},
				nsfwItem(),
				sailorItem(),
				subSubVersion(),
			}
		},
	}
}

func (a *App) submenuAdditionalSecurity() menuet.MenuItem {
	return item{
		Text: "Number and special char randomly placed",
		Clicked: func() {
			a.config.RandomPlacing = !a.config.RandomPlacing
			// if config.randomPlacing {
			// 	config.randomPlacing = false
			// } else {
			// 	config.randomPlacing = true
			// }
		},
		State: a.config.RandomPlacing}
}

func (a *App) sailorItem() menuet.MenuItem {
	if !config.Profanity.Nsfw {
		return item{
			Text: "Sailor-redneck mode (only in NSFW mode)",
		}
	}

	return item{
		Text: "Sailor-redneck mode",
		Clicked: func() {
			config.LoadDict = false
			config.Profanity.Sailor = !config.Profanity.Sailor
			if config.Profanity.Sailor {
				a.generatorVariant = SailorGenerator
				a.menu.Notification(menuet.Notification{
					Title:   "A less secure novelty setting.",
					Message: "Also using it will make you look like a juvenile asshole. Use at your own risk.",
				})
			}
			a.setMenuState()
		},
		State: a.config.Profanity.Sailor,
	}
}

func (a *App) nsfwItem() menuet.MenuItem {
	return item{
		Text: "NSFW",
		Clicked: func() {
			if a.config.Profanity.Nsfw {
				a.generatorVariant = SfwGenerator
				a.config.Profanity.Sfw = true
				a.config.Profanity.Nsfw = false
			} else {
				a.generatorVariant = NsfwGenerator
				a.config.Profanity.Sfw = false
				a.config.Profanity.Nsfw = true
			}
			a.config.LoadDict = false
			a.config.Profanity.Sailor = false
			a.setMenuState()
		},
		State: config.Profanity.Nsfw,
	}
}

func (a *App) subSubVersion() item {
	return item{
		Text:     Version,
		FontSize: 7,
	}
}
