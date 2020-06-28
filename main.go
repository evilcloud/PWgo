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
	sailorState   bool
	passLenght    int
)

type item = menuet.MenuItem

func menuItems() []item {
	// var passLenght int
	if len(adjectives) < 1 {
		adjectives = openFile("Resources/adjectives.txt")
		nouns = openFile("Resources/nouns.txt")
		fmt.Println("initial")
		passLenght := 40
		fmt.Println(passLenght)

	}

	if !loadDict && sfwDict {
		loadDict = true
		adjectives = openFile("Resources/adjectives.txt")
		nouns = openFile("Resources/nouns.txt")
		fmt.Println("SFW")
	}

	if !loadDict && nsfwDict {
		loadDict = true
		bad := openFile("Resources/bad.txt")
		if sailorRedneck {
			adjectives = bad
			nouns = bad
			fmt.Println("Hello, sailor!")
		} else {
			adjectives = append(openFile("Resources/adjectives.txt"), bad...)
			nouns = append(openFile("Resources/nouns.txt"), bad...)
			fmt.Println("NSFW")
		}
	}
	passData := append(adjectives, nouns...)
	username := strings.Title(pickRandomWord(passData)) + strings.Title(pickRandomWord(nouns))
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
					{Text: "Lenght"},
					{
						Text: "10",
						Clicked: func() {
							passLenght := 10
							fmt.Println(passLenght)
						},
					}, {
						Text: "20",
						Clicked: func() {
							passLenght := 20
							fmt.Println(passLenght)

						},
					},
					{
						Text: "30",
						Clicked: func() {
							passLenght := 30
							fmt.Println(passLenght)

						},
					},
					{
						Text: "40",
						Clicked: func() {
							passLenght := 40
							fmt.Println(passLenght)

						},
					},
					spacer,
					{Text: "NSFW",
						Clicked: func() {
							loadDict = false
							if nsfwDict {
								nsfwDict = false
								sailorState = false
								sfwDict = true
							} else {
								nsfwDict = true
								sfwDict = false
								sailorState = true
							}
						},
						State: nsfwDict},
					{
						Text: "Sailor-redneck mode",
					},
				}
			},
			// 	Children: func() []menuet.MenuItem {
			// 		return []menuet.MenuItem{
			// 			Text: "10",
			// 			Clicked: func() {
			// 				lenght := 10
			// 			},
			// 			State: true,
			// 		},
			// 	}
		},
		// item{
		// 	Text: "NSFW",
		// 	Clicked: func() {
		// 		loadDict = false
		// 		if nsfwDict {
		// 			nsfwDict = false
		// 			sailorState = false
		// 			sfwDict = true
		// 		} else {
		// 			nsfwDict = true
		// 			sfwDict = false
		// 			sailorState = true
		// 		}
		// 	},
		// 	State: nsfwDict,
		// },
		// item{
		// 	Text: "Sailor-redneck mode",
		// },
	}
}

// item{
// 	Text: "Sailor-redneck mode",
// 	if sailorState{
// 		Clicked: func() {
// 			loadDict = false
// 			if sailorRedneck {
// 				sailorRedneck = false
// 			} else {
// 				sailorRedneck = true
// 			}
// 		}
// 	},
// State: sailorState,
// }

// func passLenght() []menuet.MenuItem {
// 	return []menuet.MenuItem{
// 		Text: "10"}
// }

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

func generatePass(passLenght int, data []string) string {
	fmt.Println(passLenght)
	var generatedPass string = ""
	var i int
	// var pLeng int = passLenght - 2
	// var attempts int = 0
	// for attempts < 1000000 {
	// 	if len(generagedPass) > pLeng {
	// 		generagedPass = ""
	// 	}
	// 	if len(generagedPass) < pLeng {
	// 		generagedPass = generagedPass + pickRandomWord(data)
	// 	}
	// 	if len(generagedPass) == pLeng {
	// 		break
	// 	}
	// 	attempts += 1
	// }
	for i > 7 {
		generatedPass += pickRandomWord(data)
		fmt.Println(len(data))
	}
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
