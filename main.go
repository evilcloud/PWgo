package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	a "github.com/evilcloud/PWgo/internal/app"
	t "github.com/evilcloud/PWgo/internal/types"
	"github.com/evilcloud/PWgo/pkg/generators"

	"github.com/caseymrm/menuet"
)

// TODO: Externalise all strings
// TODO: build a production init tool (go build -ldflags="-X 'main.Version=v1.0.0'")
// TODO: build automatic packer for create-dmg 'Password machine.app' --dmg-title='Password machine'

// type credentials struct {
// 	uname struct {
// 		value string
// 		time  time.Time
// 	}
// 	pass struct {
// 		value string
// 		time  time.Time
// 	}
// }

const (
	adjFile  = "data/adjectives.txt"
	nounFile = "data/nouns.txt"
	badFile  = "data/bad.txt"
	// devVersionString = "development version\t"
	// passShort        = 8
	// passAcceptable   = 12
	// passStandard     = 20
	// passLong         = 40
)

// var (
// 	app          menuet.Application
// 	Version      string = devVersionString
// 	config       t.Settings
// 	currCreds    credentials
// 	clickedCreds credentials
// 	adjectives   []string
// 	nouns        []string
// )

type item = menuet.MenuItem

func main() {
	isDevVersion()
	getDefaults()

	config := t.Settings{}
	// if config.Profanity.Sfw == config.Profanity.Nsfw == config.Profanity.Sailor == false {
	// 	config.Profanity.Sfw = true
	// }

	// config.Info()

	// app := menuet.App()
	// app.Children = menuItems
	// app.Name = "Password machine"
	// app.Label = "com.github.evilcloud.PWgo"
	// app.AutoUpdate.Version = Version
	// app.AutoUpdate.Repo = "evilcloud/PWgo"

	// setMenuState()

	dictAdjectives := openFile(adjFile)
	dictNouns := openFile(nounFile)
	dictBad := openFile(badFile)

	sfwGenerator := generators.NewGenerator(config, app, dictNouns, dictAdjectives)
	nsfwGenerator := generators.NewGenerator(config, app, append(dictNouns, dictBad...), append(dictAdjectives, dictBad...))
	sailorGenerator := generators.NewGenerator(config, app, dictBad, dictBad)

	logicApp := a.NewApp(config, sfwGenerator, nsfwGenerator, sailorGenerator)

	logicApp.Run()
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
		config.DevVersion = true
	}
}

//testing
// // dictionaries
// func checkDictionaries() {
// 	if config.LoadDict == false {
// 		config.loadDict = true
// 		switch {
// 		case len(adjectives) < 1:
// 			config.loadDict = true
// 			adjectives = openFile(adjFile)
// 			nouns = openFile(nounFile)
// 			debugNotification("initial dictionary load")
// 		case config.profanity.sfw:
// 			config.loadDict = true
// 			adjectives = openFile(adjFile)
// 			nouns = openFile(nounFile)
// 			log.Print("SFW")
// 		case config.profanity.nsfw:
// 			bad := openFile(badFile)
// 			if config.profanity.sailor {
// 				config.loadDict = true
// 				adjectives = bad
// 				nouns = adjectives
// 				debugNotification("Hello, sailor!")
// 			} else {
// 				config.loadDict = true
// 				adjectives = append(openFile(adjFile), bad...)
// 				nouns = append(openFile(nounFile), bad...)
// 				log.Print("NSFW")
// 			}
// 		}
// 	}
// }
