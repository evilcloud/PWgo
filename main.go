package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/dustin/go-humanize"
)

// TODO: Externalise all strings
// TODO: build a production init tool (go build -ldflags="-X 'main.Version=v1.0.0'")
// TODO: build automatic packer for create-dmg 'Password machine.app' --dmg-title='Password machine'

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

type settings struct {
	passLength    int
	RandomPlacing bool
	loadDict      bool
	devVersion    bool
	profanity     struct {
		sfw    bool
		nsfw   bool
		sailor bool
	}
}

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

var (
	app          menuet.Application
	Version      string = devVersionString
	config       settings
	currCreds    credentials
	clickedCreds credentials
	adjectives   []string
	nouns        []string
)

type item = menuet.MenuItem

func init() {
	isDevVersion()
	getDefaults()

	config := settings{}
	if config.passLength < passShort {
		config.passLength = passStandard
	}
	if config.profanity.sfw == config.profanity.nsfw == config.profanity.sailor == false {
		config.profanity.sfw = true
	}

	currCreds = credentials{}
	clickedCreds = credentials{}

	// config.Info()

	app := menuet.App()
	app.Children = menuItems
	app.Name = "Password machine"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = Version
	app.AutoUpdate.Repo = "evilcloud/PWgo"

	setMenuState()
}

func main() {
	app.RunApplication()
}

func setMenuState() {
	var image string
	if config.profanity.sailor {
		image = "sailor.pdf"
	} else if config.profanity.nsfw {
		image = "nsfw.pdf"
	} else {
		image = "pw.pdf"
	}

	menuet.App().SetMenuState(&menuet.MenuState{
		Image: image,
	})
	menuet.App().MenuChanged()
}

func submenuAdditionalSecurity() menuet.MenuItem {
	return item{
		Text: "Number and special char randomly placed",
		Clicked: func() {
			if config.RandomPlacing {
				config.RandomPlacing = false
			} else {
				config.RandomPlacing = true
			}
		},
		State: config.RandomPlacing}
}

func openFile(fileName string) []string {
	ex, err := os.Executable()
	isError(err)
	exPath := filepath.Dir(ex)
	fullPath := path.Join(exPath, fileName)

	fileContent, err := ioutil.ReadFile(fullPath)
	isError(err)
	return strings.Split(string(fileContent), "\n")
}

// Adds to Version string the time of compile and trues config.devVersion if version is not changed by idflags
func isDevVersion() {
	if Version == devVersionString {
		t := time.Now()
		Version += t.String()
		config.devVersion = true
	}
}

// dictionaries
func checkDictionaries() {
	if config.loadDict == false {
		config.loadDict = true
		switch {
		case len(adjectives) < 1:
			config.loadDict = true
			adjectives = openFile(adjFile)
			nouns = openFile(nounFile)
			debugNotification("initial dictionary load")
		case config.profanity.sfw:
			config.loadDict = true
			adjectives = openFile(adjFile)
			nouns = openFile(nounFile)
			log.Print("SFW")
		case config.profanity.nsfw:
			bad := openFile(badFile)
			if config.profanity.sailor {
				config.loadDict = true
				adjectives = bad
				nouns = adjectives
				debugNotification("Hello, sailor!")
			} else {
				config.loadDict = true
				adjectives = append(openFile(adjFile), bad...)
				nouns = append(openFile(nounFile), bad...)
				log.Print("NSFW")
			}
		}
	}
}

func humaniseDuration(start time.Time) string {
	ret := humanize.Time(start)
	debugNotification(ret)
	if ret == "a long while ago" {
		ret = ""
	}
	return ret
}

// Persistent state

// Load previous state
func getDefaults() {
	clickedCreds.uname.value = menuet.Defaults().String("uname.value")
	clickedCreds.pass.value = menuet.Defaults().String("pass.value")
	config.passLength = menuet.Defaults().Integer("passLength")
}

// Save the state
func setDefaults() {
	menuet.Defaults().SetInteger("passLength", config.passLength)
	menuet.Defaults().SetString("uname.value", clickedCreds.uname.value)
	menuet.Defaults().SetString("pass.value", clickedCreds.pass.value)
}

// func (a *settings) update() {
// 	if a.passLength < passShort {
// 		a.passLength = passStandard
// 	}

// 	menuet.Defaults().SetInteger("passLength", a.passLength)
// }

// func (a *credentials) update() {
// 	menuet.Defaults().SetString("uname.value", a.uname.value)
// 	menuet.Defaults().SetString("pass.value", a.pass.value)
// }
