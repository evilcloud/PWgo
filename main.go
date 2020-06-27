package main

import (
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
	adjectives []string
	nouns      []string
)

type item = menuet.MenuItem

func main() {

	app := menuet.App()
	app.SetMenuState(&menuet.MenuState{
		Title: "Password",
	})
	app.Children = menuItems
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.pwgo"
	app.RunApplication()
}

func menuItems() []item {
	adjectives = openFile("resources/adjectives.txt")
	nouns = openFile("resources/nouns.txt")

	username := strings.Title(pickRandomWord(adjectives)) + strings.Title(pickRandomWord(nouns))
	password := generatePass(12, adjectives)
	spacer := item{}

	return []item{
		item{Text: "Username"},
		item{Text: username,
			Clicked: func() {
				clipboard.WriteAll(username)
			}},
		spacer,
		item{Text: "Password"},
		item{
			Text: password,
			Clicked: func() {
				clipboard.WriteAll(password)
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

func generatePass(lenght int, adjectives []string) string {
	password := pickRandomWord(adjectives)
	password += pickRandomWord(nouns)
	password += pickRandomWord(strings.Split("1 2 3 4 5 6 7 8 9", " "))
	password += pickRandomWord(strings.Split("! @ # $ % & * - + = ?", " "))
	return password
}

func pickRandomWord(data []string) string {
	rand.Seed(time.Now().Unix())
	return strings.Title(data[rand.Intn(len(data))])
}

// func getIcon(path string) []byte {
// 	icon, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		log.Fatal("Icon not found")
// 	}
// 	return icon
// }
