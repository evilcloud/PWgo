package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/caseymrm/menuet"
	"gopkg.in/yaml.v2"
)

// FIXME: sort out the scopes -- too many globals?
// TODO: Add time to last-clicked passwords
// TODO: Externalise all strings

func debugNotification(text string) {
	menuet.App().Notification(menuet.Notification{
		Title:   "Debug notification",
		Message: text,
	})
}

func main() {
	// debugMode := false
	// conf = loadConfig()
	setMenuState("sfw")
	app := menuet.App()
	app.Children = menuItems
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = "v0.2"
	app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.RunApplication()
}

var (
	adjectives      []string
	nouns           []string
	sfwDict         bool
	nsfwDict        bool
	sailorRedneck   bool
	loadDict        bool
	nsRandomPlace   bool
	passLenght      int
	clickedUsername string
	clickedPassword string
	conf            Configuration
)

const (
	adjFile        = "data/adjectives.txt"
	nounFile       = "data/nouns.txt"
	badFile        = "data/bad.txt"
	configFileName = "config.yaml"
	m1             = 4
	m2             = 5
	m3             = 6
	m4             = 8
)

type item = menuet.MenuItem

// Configuration
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

func saveConfig(confFileName string, conf Configuration) error {
	bytes, err := yaml.Marshal(conf)
	if err != nil {
		isError(err)
		return err
	}
	debugPrintConf()
	return ioutil.WriteFile(confFileName, bytes, 0644)
}

func loadConfig(confFileName string) Configuration {
	var conf Configuration
	_, err := os.Stat(confFileName)
	if os.IsNotExist(err) {
		isError("No configuration file found. Creating new...")
		conf = virginConfig()
		err := saveConfig(confFileName, conf)
		isError(err)
		if err != nil {
			log.Panic(err)
		}
	}
	bytes, err := ioutil.ReadFile(confFileName)
	if err != nil {
		isError(err)
		conf = virginConfig()
	} else {
		err = yaml.Unmarshal(bytes, &conf)
		if err != nil {
			isError(err)
		}
	}
	return conf
}

// func createInitialConfig() Configuration {
// 	return Configuration{
// 		Profanity: "SFW",
// 		LastName:  "",
// 		LastPass:  "",
// 	}
// }

// end of config

func menuItems() []item {
	username, password := generateUsernamePass()

	spacer := item{}

	return []item{
		item{Text: "Username"},
		item{Text: username,
			FontWeight: menuet.WeightMedium,
			FontSize:   16,
			Clicked: func() {
				clipboard.WriteAll(username)
				clickedUsername = username
			}},
		item{Text: "Password"},
		item{
			Text:       password,
			FontWeight: menuet.WeightMedium,
			Clicked: func() {
				clipboard.WriteAll(password)
				clickedPassword = password
			},
		},
		spacer,
		submenuLastItemClicked(),
		spacer,
		item{
			Text: "Settings",
			Children: func() []menuet.MenuItem {
				return []menuet.MenuItem{
					{Text: "Length (words)"},
					{
						Text: strconv.Itoa(m1),
						Clicked: func() {
							passLenght = m1
						},
						State: passLenght == m1,
					}, {
						Text: strconv.Itoa(m2),
						Clicked: func() {
							passLenght = m2
						},
						State: passLenght == m2,
					},
					{
						Text: strconv.Itoa(m3),
						Clicked: func() {
							passLenght = m3
						},
						State: passLenght == m3,
					},
					{
						Text: strconv.Itoa(m4),
						Clicked: func() {
							passLenght = m4
						},
						State: passLenght == m4,
					},
					{Text: "Additional security"},
					submenuAdditionalSecurity(),
					spacer,
					nsfwItem(),
					sailorItem(),
				}
			},
		},
	}
}

// changing the menu bar icon depending on the Profanity of profanity
func setMenuState(swearState string) {
	var image string
	switch swearState {
	case "nsfw":
		image = "nsfw.pdf"
	case "sailor":
		image = "sailor.pdf"
	default:
		image = "pw.pdf"
	}
	menuet.App().SetMenuState(&menuet.MenuState{
		Image: image,
	})
	menuet.App().MenuChanged()
}

// menu items
func sailorItem() menuet.MenuItem {
	sailorItem := item{Text: "Sailor-redneck mode (only in NSFW mode)"}
	if nsfwDict {
		sailorItem = item{Text: "Sailor-redneck mode",
			Clicked: func() {
				loadDict = false
				if sailorRedneck {
					sailorRedneck = false
					nsfwDict = true
					setMenuState("nsfw")
				} else {
					sailorRedneck = true
					setMenuState("sailor")
					menuet.App().Notification(menuet.Notification{
						Title:   "A less secure novelty setting.",
						Message: "Also using it will make you look like a juvenile asshole. Use at your own risk.",
					})
				}
			},
			State: sailorRedneck,
		}
	}
	return sailorItem
}

func nsfwItem() menuet.MenuItem {
	return item{
		Text: "NSFW",
		Clicked: func() {
			loadDict = false
			if nsfwDict {
				nsfwDict = false
				sailorRedneck = false
				sfwDict = true
				setMenuState("sfw")
			} else {
				nsfwDict = true
				sfwDict = false
				sailorRedneck = false
				setMenuState("nsfw")
			}
		},
		State: nsfwDict}
}

func submenuLastItemClicked() menuet.MenuItem {
	return item{
		Text: "Last copy-clicked",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				item{Text: "Username"},
				item{Text: clickedUsername,
					FontWeight: menuet.WeightMedium,
					FontSize:   14,
					Clicked: func() {
						clipboard.WriteAll(clickedUsername)
					}},
				item{},
				item{Text: "Password"},
				item{Text: clickedPassword,
					FontWeight: menuet.WeightMedium,
					FontSize:   14,
					Clicked: func() {
						clipboard.WriteAll(clickedPassword)
					}},
			}
		},
	}
}

func submenuAdditionalSecurity() menuet.MenuItem {
	return item{
		Text: "Number and special char randomly placed",
		Clicked: func() {
			if nsRandomPlace {
				nsRandomPlace = false
			} else {
				nsRandomPlace = true
			}
		},
		State: nsRandomPlace}
}

// general functions
func generatePass(passData []string) string {
	// FIXME: why is empty passData array's length is 2?
	if len(passData) == 2 {
		return ""
	}

	var generatedPass string = ""
	if passLenght < m1 {
		passLenght = m3
	}

	numberPosition := passLenght
	charPosition := passLenght

	number := pickRandomWord(strings.Split("1 2 3 4 5 6 7 8 9", " "))
	specialChar := pickRandomWord(strings.Split("! @ # $ % & * - + = ?", " "))
	if nsRandomPlace {
		numberPosition = pickNumberRange(passLenght)
		charPosition = pickNumberRange(passLenght)
	}

	for i := 1; i < passLenght+1; i++ {
		generatedPass += pickRandomWord(passData)
		if i == numberPosition {
			generatedPass += number
		}
		if i == charPosition {
			generatedPass += specialChar
		}
	}
	return generatedPass
}

func pickNumberRange(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

func pickRandomWord(data []string) string {
	rand.Seed(time.Now().UnixNano())
	return strings.Title(data[rand.Intn(len(data))])
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

func isError(err error) {
	if err != nil {
		log.Println(err)
		menuet.App().Notification(menuet.Notification{
			Title:        "Error!",
			Message:      err.Error(),
			ActionButton: "Quit app",
			CloseButton:  "Close notification",
		})
	}
}

func generateUsernamePass() (string, string) {
	if loadDict == false {
		loadDict = true
		switch {
		case len(adjectives) < 1:
			loadDict = true
			adjectives = openFile(adjFile)
			nouns = openFile(nounFile)
			log.Println("initial dictionary load")
		case sfwDict:
			loadDict = true
			adjectives = openFile(adjFile)
			nouns = openFile(nounFile)
			log.Print("SFW")
		case nsfwDict:
			bad := openFile(badFile)
			if sailorRedneck {
				loadDict = true
				adjectives = bad
				nouns = adjectives
				log.Println("Hello, sailor!")
			} else {
				loadDict = true
				adjectives = append(openFile(adjFile), bad...)
				nouns = append(openFile(nounFile), bad...)
				log.Print("NSFW")
			}
		}
	}

	passData := append(adjectives, nouns...)
	usernameUncleaned := strings.Title(pickRandomWord(adjectives)) + strings.Title(pickRandomWord(nouns))
	reg, err := regexp.Compile("[^a-zA-Z]+")
	isError(err)
	username := reg.ReplaceAllString(usernameUncleaned, "")
	if len(passData) == 2 {
		username = "⚠️ NO DICTIONARY FOUND"
	}
	password := generatePass(passData)
	// FIXED: add proper adjectives and nouns to sailor password -- NO. Let these juvies suffer from inferior protection
	clipboard.WriteAll(password)

	return username, password
}
