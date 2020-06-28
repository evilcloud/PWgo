package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

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
	passLenght    int
)

type item = menuet.MenuItem

func menuItems() []item {
	adjFile := "Resources/adjectives.txt"
	nounFile := "Resources/nouns.txt"
	badFile := "Resources/bad.txt"
	if len(adjectives) < 1 {
		adjectives = openFile(adjFile)
		nouns = openFile(nounFile)
		fmt.Println("initial")
		passLenght := 40
		fmt.Println(passLenght)

	}

	if !loadDict && sfwDict {
		loadDict = true
		adjectives = openFile(adjFile)
		nouns = openFile(nounFile)
		fmt.Println("SFW")
	}

	if !loadDict && nsfwDict {
		loadDict = true
		bad := openFile(badFile)
		adjectives = append(openFile(adjFile), bad...)
		nouns = append(openFile(nounFile), bad...)
		fmt.Println("NSFW")
	}

	if sailorRedneck && !loadDict {
		loadDict = true
		bad := openFile(badFile)
		adjectives = bad
		nouns = bad
		fmt.Println("Hello, sailor!")
	}

	passData := append(adjectives, nouns...)
	username := strings.Title(pickRandomWord(adjectives)) + strings.Title(pickRandomWord(nouns))
	fmt.Println(len(adjectives), len(nouns), len(passData))
	password := generatePass(passLenght, passData)
	spacer := item{}

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
						Text: "4",
						Clicked: func() {
							passLenght := 4
							fmt.Println(passLenght)
						},
					}, {
						Text: "5",
						Clicked: func() {
							passLenght := 5
							fmt.Println(passLenght)

						},
					},
					{
						Text: "6",
						Clicked: func() {
							passLenght := 6
							fmt.Println(passLenght)

						},
					},
					{
						Text: "7",
						Clicked: func() {
							passLenght := 7
							fmt.Println(passLenght)

						},
					},
					spacer,
					{Text: "NSFW",
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
					{
						Text: "Sailor-redneck mode",
						Clicked: func() {
							loadDict = false
							if sailorRedneck {
								sailorRedneck = false
								nsfwDict = true
							} else {
								sailorRedneck = true
								nsfwDict = false
							}
						},
						State: sailorRedneck,
					},
				}
			},
		},
	}
}

func openFile(fileName string) []string {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("File does not exist")
		} else {
			log.Fatal(err)
		}
	}
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(fileContent), "\n")
}

func generatePass(passLenght int, passData []string) string {
	fmt.Println(len(passData))
	var generatedPass string = ""
	// for i > passLenght {
	generatedPass += pickRandomWord(passData)
	generatedPass += pickRandomWord(passData)
	generatedPass += pickRandomWord(passData)
	generatedPass += pickRandomWord(passData)
	generatedPass += pickRandomWord(passData)
	// }
	generatedPass += pickRandomWord(strings.Split("1 2 3 4 5 6 7 8 9", " "))
	generatedPass += pickRandomWord(strings.Split("! @ # $ % & * - + = ?", " "))
	fmt.Println(generatedPass)
	return generatedPass
}

func pickRandomWord(data []string) string {
	rand.Seed(time.Now().UnixNano())
	return strings.Title(data[rand.Intn(len(data))])
}

func main() {
	app := menuet.App()
	app.SetMenuState(&menuet.MenuState{
		Title: "PWgo",
		// Image: "22",
	})
	app.Children = menuItems
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = "v0.1"
	app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.RunApplication()
}
