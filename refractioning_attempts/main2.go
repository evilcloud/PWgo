package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/caseymrm/menuet"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	TimeCreated      time.Time `yaml: "TimeCreated`
	Profanity        string    `yaml: "Profanity"`
	LastName         string    `yaml: "LastName`
	LastPass         string    `yaml: "LastPass`
	PassLenWords     int       `yaml: "PassLenWords"`
	PassLenChars     int       `yaml: "PassLenChars`
	PassLenCharsPick bool      `yaml: "PassLenCharsPick`
	RandomPlacement  bool      `yaml: "RandomPlacement`
	LengthsWords     []int     `yaml: "lengthsWords`
	LengthsChars     []int     `yaml: "lengthsChars`
}

func virginConfig() Configuration {
	return Configuration{
		TimeCreated:      time.Now(),
		Profanity:        "sfw",
		LastName:         "",
		LastPass:         "",
		PassLenWords:     6,
		PassLenChars:     20,
		PassLenCharsPick: true,
		RandomPlacement:  false,
		LengthsWords:     []int{4, 6, 8, 12},
		LengthsChars:     []int{8, 12, 20, 40},
	}
}

// type passLenghtWorlds struct {
// }

type item = menuet.MenuItem

const (
	gui                 = false
	confFileName string = "config.yaml"
)

var config Configuration

// var (
// 	conf Configuration
// )

func main() {
	config = loadConfig(confFileName)
	// log.Println(config)
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: config.Profanity,
	})
	app := menuet.App()
	app.Name = "IDmixer"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.AutoUpdate.Version = "v1.0"
	app.Children = menuItems
	app.RunApplication()
}

// Utilities
func popupIfErr(data string) {
	if data != "" {
		switch gui {
		case true:
			log.Println(data)
		case false:
			log.Println("that's for GUI")
		}
	}
}

func menuItems() []item {
	adjectives := openFile("data/adjectives.txt")
	nouns := openFile("data/nouns.txt")
	username := generateUsername(adjectives, nouns)
	password := "temporary password"

	return []item{
		item{Text: "Username"},
		item{Text: username,
			FontWeight: menuet.WeightMedium,
			FontSize:   16,
			Clicked: func() {
				clipboard.WriteAll(username)
				config.LastName = username
				saveConfig(confFileName, config)
			}},
		item{},
		item{Text: "Password"},
		item{Text: password,
			FontWeight: menuet.WeightMedium,
			FontSize:   16,
			Clicked: func() {
				clipboard.WriteAll(password)
				config.LastPass = password
				saveConfig(confFileName, config)
			}},
		item{},
		item{},
		// submenuSettings(),
		submenuLastItemClicked(),
		item{},
		submenuSettings(),
	}
}

func submenuSettings() item {
	return item{
		Text: "Settings",
		Children: func() []item {
			return []item{
				submenuLengthPassWords(),
				item{},
				submenuAdditionalSecurity(),
				item{},
				// nsfwItem(),
				// sailorItem(),
			}
		},
	}
}

func submenuLengthPassWords() item {
	// o1 := 8
	// o2 := 12
	// o3 := 20
	// o4 := 40

	return item{
		Text: "Length (words)",
		// {
		// 	Text: strconv.Itoa(m1),
		// 	Clicked: func() {
		// 		passLenght = m1
		// 	},
		// 	State: passLenght == m1,
		// }, {
		// 	Text: strconv.Itoa(m2),
		// 	Clicked: func() {
		// 		passLenght = m2
		// 	},
		// 	State: passLenght == m2,
		// },
		// {
		// 	Text: strconv.Itoa(m3),
		// 	Clicked: func() {
		// 		passLenght = m3
		// 	},
		// 	State: passLenght == m3,
		// },
		// {
		// 	Text: strconv.Itoa(m4),
		// 	Clicked: func() {
		// 		passLenght = m4
		// 	},
		// 	State: passLenght == m4,
		// },
		// {Text: "Additional security"},
	}
}

func submenuAdditionalSecurity() item {
	return item{
		Text: "Additional security",
		Clicked: func() {
			if config.RandomPlacement {
				config.RandomPlacement = false
				// saveConfig(confFileName, config)

			} else {
				config.RandomPlacement = true
				// saveConfig(confFileName, config)
			}
			saveConfig(confFileName, config)
		},
		State: config.RandomPlacement,
	}
}

func submenuLastItemClicked() item {
	return item{
		Text: "Last copy-clicked",
		Children: func() []item {
			return []item{
				item{Text: "Username"},
				item{Text: config.LastName,
					FontWeight: menuet.WeightMedium,
					FontSize:   14,
					Clicked: func() {
						clipboard.WriteAll(config.LastName)
					}},
				item{},
				item{Text: "Password"},
				item{Text: config.LastPass,
					FontWeight: menuet.WeightMedium,
					FontSize:   14,
					Clicked: func() {
						clipboard.WriteAll(config.LastPass)
					}},
			}
		},
	}
}

// Config

func loadConfig(confFileName string) Configuration {
	var conf Configuration
	_, err := os.Stat(confFileName)
	if os.IsNotExist(err) {
		popupIfErr("No configuration file found. Creating new...")
		conf = virginConfig()
		err := saveConfig(confFileName, conf)
		popupIfErr(err.Error())
		if err != nil {
			log.Panic(err)
		}
	}
	bytes, err := ioutil.ReadFile(confFileName)
	if err != nil {
		popupIfErr(err.Error())
		conf = virginConfig()
	} else {
		err = yaml.Unmarshal(bytes, &conf)
		if err != nil {
			popupIfErr(err.Error())
		}
	}
	return conf
}

func saveConfig(confFileName string, conf Configuration) error {
	bytes, err := yaml.Marshal(conf)
	if err != nil {
		popupIfErr(err.Error())
		return err
	}
	debugPrintConf()
	return ioutil.WriteFile(confFileName, bytes, 0644)
}

func debugPrintConf() {
	p, _ := json.MarshalIndent(config, "", "	")
	log.Println(p)
}

// generators

func generateUsername(adjectives []string, nouns []string) string {
	adjective := pickRandomWord(adjectives)
	noun := pickRandomWord(nouns)
	return adjective + noun
}

func pickRandomWord(words []string) string {
	rand.Seed(time.Now().UnixNano())
	return strings.Title(words[rand.Intn(len(words))])
}

// data loaders
func openFile(fileName string) []string {
	ex, err := os.Executable()
	popupIfErr(err.Error())
	exPath := filepath.Dir(ex)
	fullPath := path.Join(exPath, fileName)

	fileContent, err := ioutil.ReadFile(fullPath)
	popupIfErr(err.Error())
	return strings.Split(string(fileContent), "\n")
}
