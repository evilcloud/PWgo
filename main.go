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
	adjectives []string
	nouns      []string
)

type item = menuet.MenuItem

func main() {
	app := menuet.App()
	app.SetMenuState(&menuet.MenuState{
		Title: "PWgo",
	})
	app.Children = menuItems
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = "v0.1"
	app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.RunApplication()
}

func menuItems() []item {
	if len(adjectives) < 1 {
		adjectives = openFile("Resources/adjectives.txt")
		nouns = openFile("Resources/nouns.txt")
		fmt.Println("ok")
	}

	username := strings.Title(pickRandomWord(adjectives)) + strings.Title(pickRandomWord(nouns))
	password := generatePass(12, adjectives)
	spacer := item{}

	return []item{
		item{Text: "Username"},
		item{Text: username,
			FontWeight: menuet.WeightMedium,
			Clicked: func() {
				clipboard.WriteAll(username)
			}},
		spacer,
		item{Text: "Password"},
		item{
			Text:       password,
			FontWeight: menuet.WeightMedium,
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
