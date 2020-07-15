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
	"github.com/dustin/go-humanize"
)

// FIXME: sort out the scopes -- too many globals?
// TODO: Externalise all strings
// FIXME: humanise option doesn't work properly

func debugNotification(text string) {
	// log.Println(text)
	// menuet.App().Notification(menuet.Notification{
	// 	Title:   "Debug notification",
	// 	Message: text,
	// })
}

type credentials struct {
	uname struct {
		value string
		time  time.Time
	}
	pass struct {
		value string
		time  time.Time
	}
}

type settings struct {
	passLength    int
	RandomPlacing bool
	loadDict      bool
	profanity     struct {
		sfw    bool
		nsfw   bool
		sailor bool
	}
}

func (obj *settings) Info() {
	if obj.passLength < passShort {
		obj.passLength = passStandard
	}
	if obj.profanity.sfw == obj.profanity.nsfw == obj.profanity.sailor == false {
		obj.profanity.sfw = true
	}
}

const (
	adjFile        = "data/adjectives.txt"
	nounFile       = "data/nouns.txt"
	badFile        = "data/bad.txt"
	passShort      = 4
	passAcceptable = 5
	passStandard   = 6
	passLong       = 8
)

var (
	config       settings
	currCreds    credentials
	clickedCreds credentials
	adjectives   []string
	nouns        []string
)

type item = menuet.MenuItem

func main() {
	config := settings{}
	config.Info()

	currCreds = credentials{}
	clickedCreds = credentials{}

	setMenuState()

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
				clickedCreds.uname.value = details
				clickedCreds.uname.time = time.Now()
			case "password":
				clickedCreds.pass.value = details
				clickedCreds.pass.time = time.Now()
			}
		},
	}
}

func menuItems() []item {
	currCreds.uname.value, currCreds.pass.value = generateUsernamePass()

	spacer := item{}

	return []item{
		item{Text: "Username"},
		menuDisplayCredential(currCreds.uname.value, "username"),
		item{Text: "Password"},
		menuDisplayCredential(currCreds.pass.value, "password"),
		spacer,
		submenuLastItemClicked(),
		spacer,
		item{
			Text: "Words of Wisdom"},
		wordsOfWisdom(),
		spacer,
		item{
			Text: "Settings",
			Children: func() []menuet.MenuItem {
				return []menuet.MenuItem{
					{Text: "Length (words)"},
					{
						Text: strconv.Itoa(passShort),
						Clicked: func() {
							config.passLength = passShort
						},
						State: config.passLength == passShort,
					}, {
						Text: strconv.Itoa(passAcceptable),
						Clicked: func() {
							config.passLength = passAcceptable
						},
						State: config.passLength == passAcceptable,
					},
					{
						Text: strconv.Itoa(passStandard),
						Clicked: func() {
							config.passLength = passStandard
						},
						State: config.passLength == passStandard,
					},
					{
						Text: strconv.Itoa(passLong),
						Clicked: func() {
							config.passLength = passLong
						},
						State: config.passLength == passLong,
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

func setMenuState() {
	var image string

	if config.profanity.sailor {
		image = "sailor.pdf"
	} else if config.profanity.nsfw {
		image = "nsfw.pdf"
	} else {
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
	if config.profanity.nsfw {
		sailorItem = item{Text: "Sailor-redneck mode",
			Clicked: func() {
				config.loadDict = false
				if config.profanity.sailor {
					config.profanity.sailor = false
					config.profanity.nsfw = true
					// setMenuState()
				} else {
					config.profanity.sailor = true
					// setMenuState()
					menuet.App().Notification(menuet.Notification{
						Title:   "A less secure novelty setting.",
						Message: "Also using it will make you look like a juvenile asshole. Use at your own risk.",
					})
				}
				setMenuState()
			},
			State: config.profanity.sailor,
		}
	}
	return sailorItem
}

func nsfwItem() menuet.MenuItem {
	return item{
		Text: "NSFW",
		Clicked: func() {
			config.loadDict = false
			if config.profanity.nsfw {
				config.profanity.sfw = true
				config.profanity.nsfw = false
			} else {
				config.profanity.sfw = false
				config.profanity.nsfw = true
			}
			config.profanity.sailor = false
			setMenuState()
		},
		State: config.profanity.nsfw}
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
				menuDisplayCredential(clickedCreds.uname.value, "clicked"),
				item{},
				item{Text: "Password"},
				item{
					Text:     humaniseDuration(clickedCreds.pass.time),
					FontSize: 10,
				},
				menuDisplayCredential(clickedCreds.pass.value, "clicked"),
			}
		},
	}
}

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
	if config.passLength < passShort {
		config.passLength = passStandard
	}

	numberPosition := config.passLength
	charPosition := config.passLength

	number := pickRandomWord(strings.Split("1 2 3 4 5 6 7 8 9", " "))
	specialChar := pickRandomWord(strings.Split("! @ # $ % & * - + = ?", " "))
	if config.RandomPlacing {
		numberPosition = pickNumberRange(config.passLength)
		charPosition = pickNumberRange(config.passLength)
	}

	for i := 1; i < config.passLength+1; i++ {
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
	if config.loadDict == false {
		config.loadDict = true
		switch {
		case len(adjectives) < 1:
			config.loadDict = true
			adjectives = openFile(adjFile)
			nouns = openFile(nounFile)
			log.Println("initial dictionary load")
		case config.profanity.sfw:
			config.loadDict = true
			adjectives = openFile(adjFile)
			nouns = openFile(nounFile)
			log.Print("SFW")
		case config.profanity.nsfw:
			bad := openFile(badFile)
			if config.profanity.sailor {
				config.loadDict = true
				adjectives = bad
				nouns = adjectives
				log.Println("Hello, sailor!")
			} else {
				config.loadDict = true
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
	return humanize.Time(start)
}
