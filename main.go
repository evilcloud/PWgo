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
	randomPlacing bool
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
		// t := time.Now()
		Version += executableLastMod()
		config.devVersion = true
	}
}

//testing
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
