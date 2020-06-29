package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"path/filepath"
	"path"

	"github.com/atotto/clipboard"
	"github.com/caseymrm/menuet"
)

var (
	adjectives    []string
	nouns         []string
	sfwDict       bool
	nsfwDict      bool
	sailorRedneck bool
	loadDict      bool
	nsRandomPlace bool
	passLenght    int
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

	if len(adjectives) < 1 {
		adjectives = openFile(adjFile)
		nouns = openFile(nounFile)
		log.Println("initial")
	}

	if !loadDict && sfwDict {
		loadDict = true
		adjectives = openFile(adjFile)
		nouns = openFile(nounFile)
		log.Println("SFW")
	}

	if !loadDict && nsfwDict {
		loadDict = true
		bad := openFile(badFile)
		if sailorRedneck {
			adjectives = bad
			nouns = bad
			log.Println("Hello, sailor!")
		} else {
			adjectives = append(openFile(adjFile), bad...)
			nouns = append(openFile(nounFile), bad...)
			log.Println("NSFW")
		}
	}

	passData := append(adjectives, nouns...)
	username := strings.Title(pickRandomWord(adjectives)) + strings.Title(pickRandomWord(nouns))
// TODO: remove hyphens
	password := generatePass(passData)
// TODO: add proper adjectives and nouns to sailor password
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
			}},
		item{Text: "Password"},
		item{
			Text:       password,
			FontWeight: menuet.WeightMedium,
			Clicked: func() {
				clipboard.WriteAll(password)
			},
		},
		spacer,
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

	// if _, err := os.Stat(fileName); err != nil {
	// 	if os.IsNotExist(err) {
	// 		log.Fatal("File does not exist")
	// 	} else {
	// 		log.Println(err)
	// 		// log.Fatal(err)
	// 	}
	// }

	// pwd, err := os.Getwd()
	// isError(err)
	// fullPath := path.Join(pwd, fileName)

	fileContent, err := ioutil.ReadFile(fullPath)
	isError(err)
	return strings.Split(string(fileContent), "\n")
}

func isError(err error) {
	if err != nil {
		log.Println(err)
	}
	
}
func generatePass(passData []string) string {
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

func main() {
	// log.SetOutput(ioutil.Discard)
	app := menuet.App()
	app.SetMenuState(&menuet.MenuState{
		// Title: "PWgo",
		Image: "pw.pdf",
	})
	app.Children = menuItems
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.PWgo"
	// app.AutoUpdate.Version = "v0.1"
	// app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.RunApplication()
}
