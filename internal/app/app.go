package app

import (
	"html"
	"math/rand"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/caseymrm/menuet"
	t "github.com/evilcloud/PWgo/internal/types"
)

var (
	Version      string = devVersionString
	config       t.Settings
	currCreds    credentials
	clickedCreds credentials
)

const (
	adjFile          = "data/adjectives.txt"
	nounFile         = "data/nouns.txt"
	badFile          = "data/bad.txt"
	devVersionString = "development version\t"
	passShort        = 8
	passAcceptable   = 12
	passStandard     = 20
	passLong         = 40
)

type item = menuet.MenuItem

type credentials struct {
	uname struct {
		value string
		time  time.Time
	}
	pass struct {
		value string
		time  time.Time
	}
}

type Generator interface {
	//EmojisSeeded(num int, seed time.Time) string
	Username() string
	Password() string
	WoW() string
}

type GeneratorVariant int

const (
	UndefinedGenerator = iota
	SfwGenerator
	NsfwGenerator
	SailorGenerator
)

type App struct {
	config t.Settings
	menu   *menuet.Application

	currCreds    credentials
	clickedCreds credentials

	sfwGenerator    Generator
	nsfwGenerator   Generator
	sailorGenerator Generator

	generatorVariant GeneratorVariant
}

func NewApp(config t.Settings, sfwGenerator, nsfwGenerator, sailorGenerator Generator) *App {
	a := &App{
		config:          config,
		menu:            menuet.App(),
		sfwGenerator:    sfwGenerator,
		nsfwGenerator:   nsfwGenerator,
		sailorGenerator: sailorGenerator,
	}

	a.getDefaults()

	a.menu.Name = "Password machine"
	a.menu.Label = "com.github.evilcloud.PWgo"
	a.menu.AutoUpdate.Version = Version
	a.menu.AutoUpdate.Repo = "evilcloud/PWgo"
	a.menu.Children = a.menuItems()
	return a
}

func (a *App) Run() {
	a.setMenuState()
	a.menu.RunApplication()

}

func (a *App) generatePassword() string {
	switch a.generatorVariant {
	case SfwGenerator:
		return a.sfwGenerator.Password()
	case NsfwGenerator:
		return a.nsfwGenerator.Password()
	case SailorGenerator:
		return a.sailorGenerator.Password()
	}
	return ""
}

func (a *App) generateUsername() string {
	switch a.generatorVariant {
	case SfwGenerator:
		return a.sfwGenerator.Username()
	case NsfwGenerator:
		return a.nsfwGenerator.Username()
	case SailorGenerator:
		return a.sailorGenerator.Username()
	}
	return ""
}

func (a *App) generateWoW() string {
	switch a.generatorVariant {
	case SfwGenerator:
		return a.sfwGenerator.WoW()
	case NsfwGenerator:
		return a.nsfwGenerator.WoW()
	case SailorGenerator:
		return a.sailorGenerator.WoW()
	}
	return ""
}

func (*App) emojisSeeded(num int, seed time.Time) string {
	var emojis string = ""
	for i := 0; i < num; i++ {
		rand.Seed(seed.AddDate(i, i, i).UnixNano())
		emojiNumber := strconv.Itoa((rand.Intn(64)) + 128640)
		emojis += html.UnescapeString("&#" + emojiNumber + ";")
	}
	return emojis
}

func (a *App) setMenuState() {
	var image string
	if a.config.Profanity.Sailor {
		image = "sailor.pdf"
	} else if a.config.Profanity.Nsfw {
		image = "nsfw.pdf"
	} else {
		image = "pw.pdf"
	}

	menuet.App().SetMenuState(&menuet.MenuState{
		Image: image,
	})
	menuet.App().MenuChanged()
}

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
				a.clickedCreds.uname.value = details
			case "password":
				a.clickedCreds.pass.value = details
			}
			setDefaults()
			// credentials.update()
		},
	}
}

func (a *App) passwordMachineText() string {
	if !a.config.DevVersion {
		return "Password machine"
	}
	return "Password machine\t" + a.emojisSeeded(2, executableLastModTime())
}
