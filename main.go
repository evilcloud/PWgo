package main

import (
	"fmt"
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
	"github.com/hako/durafmt"
)

// FIXME: sort out the scopes -- too many globals?
// TODO: Externalise all strings

func debugNotification(text string) {
	// log.Println(text)
	// menuet.App().Notification(menuet.Notification{
	// 	Title:   "Debug notification",
	// 	Message: text,
	// })
}

type credentials struct {
	Username string
	Password string
	uname    struct {
		value string
		time  time.Time
	}
	pass struct {
		value string
		time  time.Time
	}

	time struct {
		last time.Time
	}
}

type settings struct {
	PassLenght    int
	Profanity     string
	RandomPlacing bool
}

func (obj *settings) Info() {
	if obj.PassLenght == 0 {
		obj.PassLenght = 6
	}
	if obj.Profanity == "" {
		obj.Profanity = "sfw"
	}
}

const (
	adjFile  = "data/adjectives.txt"
	nounFile = "data/nouns.txt"
	badFile  = "data/bad.txt"
	m1       = 4
	m2       = 5
	m3       = 6
	m4       = 8
)

var (
	config        settings
	currCreds     credentials
	clickedCreds  credentials
	adjectives    []string
	nouns         []string
	sfwDict       bool
	nsfwDict      bool
	sailorRedneck bool
	loadDict      bool
	passLenght    int
)

type item = menuet.MenuItem

func main() {
	// debugMode := false
	config := settings{}
	config.Info()

	currCreds = credentials{}
	clickedCreds = credentials{}

	setMenuState(config.Profanity)

	app := menuet.App()
	app.Children = menuItems
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = "v0.2"
	app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.RunApplication()
}

func menuDisplayCredential(details string, mode string) item {
	return item{
		Text:       details,
		FontWeight: menuet.WeightMedium,
		FontSize:   16,
		Clicked: func() {
			clipboard.WriteAll(details)
			switch mode {
			case "username":
				clickedCreds.Username = details
				clickedCreds.uname.time = time.Now()
			case "password":
				clickedCreds.Password = details
				clickedCreds.pass.time = time.Now()
			}
		},
	}
}

func menuItems() []item {
	currCreds.Username, currCreds.Password = generateUsernamePass()

	spacer := item{}

	return []item{
		item{Text: "Username"},
		menuDisplayCredential(currCreds.Username, "username"),
		item{Text: "Password"},
		menuDisplayCredential(currCreds.Password, "password"),
		spacer,
		submenuLastItemClicked(),
		spacer,
		item{
			Text: "Worlds of Wisdom"},
		wordsOfWisdom(),
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

func setMenuState(profanity string) {
	var image string
	debugNotification("setMenuState: " + config.Profanity)
	switch profanity {
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
				item{
					Text:     humaniseDuration(clickedCreds.uname.time),
					FontSize: 10,
				},
				menuDisplayCredential(clickedCreds.Username, "clicked"),
				item{},
				item{Text: "Password"},
				item{
					Text:     humaniseDuration(clickedCreds.pass.time),
					FontSize: 10,
				},
				menuDisplayCredential(clickedCreds.Password, "clicked"),
			}
		},
	}
}

// func whatTimeAgo(start time.Time) string {
// 	end := time.Now()
// 	return humaniseDuration(start, end)
// }

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

func wordsOfWisdom() item {
	word1 := pickRandomWord(adjectives)
	word2 := pickRandomWord(nouns)
	return menuDisplayCredential(word1+" "+word2, "none")
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
	if config.RandomPlacing {
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

func humaniseDuration(start time.Time) string {
	end := time.Now()
	diff := end.Sub(start)
	var empty time.Time
	if start == empty {
		return ""
	}
	duration, err := durafmt.ParseString(diff.String())
	if err != nil {
		fmt.Println(err)
	}
	dur := duration.Duration().Round(time.Second)
	return dur.String() + " ago"
}
