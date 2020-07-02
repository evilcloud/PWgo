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
)

// FIX: sort out the scopes -- too many globals?

func main() {
	debugMode := false

	app := menuet.App()

	if debugMode {
		log.Println("Debugging mode ON")
		app.SetMenuState(&menuet.MenuState{Title: "PWgo"})
	} else {
		log.SetOutput(ioutil.Discard)
		app.SetMenuState(&menuet.MenuState{Image: "pw.pdf"})

	}

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
)

type item = menuet.MenuItem

func menuItems() []item {
	const adjFile = "data/adjectives.txt"
	const nounFile = "data/nouns.txt"
	const badFile = "data/bad.txt"
	const m1 = 4
	const m2 = 5
	const m3 = 6
	const m4 = 8

	// BUG: NSFW works only once
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
	spacer := item{}

	sailorItem := item{Text: "Sailor-redneck mode (only in NSFW mode)"}
	if nsfwDict {
		sailorItem = item{Text: "Sailor-redneck mode",
			Clicked: func() {
				loadDict = false
				if sailorRedneck {
					sailorRedneck = false
					nsfwDict = true
				} else {
					sailorRedneck = true
					menuet.App().Notification(menuet.Notification{
						Title:   "A less secure novelty setting.",
						Message: "Also using it will make you look like a juvenile asshole.Use at your own risk, kiddo.",
					})
				}
			},
			State: sailorRedneck,
		}
	}

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
		item{
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
					spacer,
					item{Text: "Password"},
					item{Text: clickedPassword,
						FontWeight: menuet.WeightMedium,
						FontSize:   14,
						Clicked: func() {
							clipboard.WriteAll(clickedPassword)
						}},
				}
			},
		},
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
					spacer,
					{Text: "Additional security"},
					{Text: "Number and special char randomly placed",
						Clicked: func() {
							if nsRandomPlace {
								nsRandomPlace = false
							} else {
								nsRandomPlace = true
							}
						},
						State: nsRandomPlace},
					spacer,
					{
						Text: "NSFW",
						Clicked: func() {
							loadDict = false
							if nsfwDict {
								nsfwDict = false
								sailorRedneck = false
								sfwDict = true
							} else {
								nsfwDict = true
								sfwDict = false
								sailorRedneck = false
							}
						},
						State: nsfwDict},
					sailorItem,
				}
			},
		},
	}
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

func generatePass(passData []string) string {

	// ISSUE: why is empty passData array's length is 2?
	if len(passData) == 2 {
		return ""
	}

	var generatedPass string = ""
	if passLenght < 4 {
		passLenght = 6
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
