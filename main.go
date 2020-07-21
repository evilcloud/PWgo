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

// TODO: Externalise all strings

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
	passShort      = 8
	passAcceptable = 12
	passStandard   = 20
	passLong       = 40
	appVersion     = "v1.0.1\t0721.21A"
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
	app.Name = "Password machine"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = appVersion
	app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.RunApplication()
}

// menu items
func menuItems() []item {
	checkDictionaries()
	currCreds.uname.value = generateUsername()
	currCreds.pass.value = generatePassword()
	clipboard.WriteAll(currCreds.pass.value)
	wow := generateWoW()

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
		menuDisplayCredential(wow, "none"),
		spacer,
		submenuSettings(),
	}
}

func menuDisplayCredential(details, mode string) item {
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

func submenuSettings() item {
	subSubLengthItem := func(length int) item {
		return item{
			Text: strconv.Itoa(length),
			Clicked: func() {
				config.passLength = length
			},
			State: config.passLength == length,
		}
	}
	return item{
		Text: "Settings",
		Children: func() []menuet.MenuItem {
			return []menuet.MenuItem{
				{Text: "Password length"},
				subSubLengthItem(passShort),
				subSubLengthItem(passAcceptable),
				subSubLengthItem(passStandard),
				subSubLengthItem(passLong),
				item{},
				// {Text: "Additional security"},
				// submenuAdditionalSecurity(),
				// item{},
				item{Text: "Level of profanity"},
				nsfwItem(),
				sailorItem(),
				subSubAppVersion(),
			}
		},
	}
}

func sailorItem() menuet.MenuItem {
	sailorItem := item{Text: "Sailor-redneck mode (only in NSFW mode)"}
	if config.profanity.nsfw {
		sailorItem = item{Text: "Sailor-redneck mode",
			Clicked: func() {
				config.loadDict = false
				if config.profanity.sailor {
					config.profanity.sailor = false
					config.profanity.nsfw = true
				} else {
					config.profanity.sailor = true
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

func subSubAppVersion() item {
	return item{
		Text:     appVersion,
		FontSize: 7,
	}
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

// general functions
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

// Messaging
func isError(err error) {
	if err != nil {
		log.Println(err)
		popupMessage("Error", err.Error())
	}
}

// Creates a banner with the predetermined title and message. Dependency on menuet library
func popupMessage(title, message string) {
	menuet.App().Notification(menuet.Notification{
		Title:   title,
		Message: message,
	})
}

// Generators

// Returns a joined pair of words from generateAdjectiveNounPair() function
func generateUsername() string {
	return strings.Join(generateAdjectiveNounPair(), "")
}

// Returns a joined with a space pair of word from generateAdjectiveNounPair() function
func generateWoW() string {
	return strings.Join(generateAdjectiveNounPair(), " ")
}

// Generates a pair of random adjective and noun (only alpha characters) and returns as an array
// Used mostly by Words of Wisdom and Username functions
func generateAdjectiveNounPair() []string {
	var ret []string
	reg, err := regexp.Compile("[^a-zA-Z]+")
	isError(err)

	adj := reg.ReplaceAllLiteralString(strings.Title(pickRandomWord(adjectives)), "")
	nou := reg.ReplaceAllLiteralString(strings.Title(pickRandomWord(nouns)), "")
	ret = append(ret, []string{adj, nou}...)
	return ret
}

// Generates a string of random words from adjectives and nouns dictionaries, and returns as a string
func generatePassword() string {
	// FIXME: put the check in the init part of the code
	if config.passLength < passShort {
		config.passLength = passStandard
	}

	var pass []string
	totalDict := append(adjectives, nouns...)
	for i := 1; i < 1000; i++ {
		pass = append([]string{pickRandomWord(totalDict)}, pass...)
		lenPassAlpha := len(strings.Join(pass, "")) - 2
		if lenPassAlpha == config.passLength {
			pass := insertRandomNumChar(pass)
			return strings.Join(pass, "")
		} else if lenPassAlpha > config.passLength {
			pass = nil
		}
	}
	popupMessage("Password failure", "Failed to match pass to length!")
	return "Failed to match pass to length"
}

// Generates random number (except 0) and some special character (from the list). Returns as an array
func insertRandomNumChar(data []string) []string {
	numeral := pickRandomWord(strings.Split("1 2 3 4 5 6 7 8 9", " "))
	char := pickRandomWord(strings.Split("! @ # $ % & * - + = ?", " "))
	// lenData := len(data)
	// var numPosition, charPosition int

	// if config.RandomPlacing {
	// 	numPosition := pickNumberRange(lenData)
	// 	charPosition := pickNumberRange(lenData + 1)
	// } else {
	// 	numPosition := lenData
	// 	charPosition := lenData + 1
	// }
	data = append(data, []string{numeral}...)
	data = append(data, []string{char}...)
	return data
}

// dictionaries
func checkDictionaries() {
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
}

func humaniseDuration(start time.Time) string {
	ret := humanize.Time(start)
	log.Println(ret)
	if ret == "a long while ago" {
		ret = ""
	}
	return ret
}
