package main

import (
	"html"
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
// TODO: build a production init tool (go build -ldflags="-X 'main.Version=v1.0.0'")
// TODO: build automatic packer for create-dmg 'Password machine.app' --dmg-title='Password machine'

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
	devVersion    bool
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
	adjFile          = "data/adjectives.txt"
	nounFile         = "data/nouns.txt"
	badFile          = "data/bad.txt"
	devVersionString = "development version\t"
	passShort        = 8
	passAcceptable   = 12
	passStandard     = 20
	passLong         = 40
)

var (
	Version      string = devVersionString
	config       settings
	currCreds    credentials
	clickedCreds credentials
	adjectives   []string
	nouns        []string
)

type item = menuet.MenuItem

func getDefaults() {
	clickedCreds.uname.value = menuet.Defaults().String("uname.value")
	// log.Println(menuet.Defaults().String("uname.time"))
	clickedCreds.pass.value = menuet.Defaults().String("pass.value")
	// clickedCreds.pass.time, _ = time.Parse("2020-07-22 15:33:17.865231 +0800 HKT m=+16.132953837", menuet.Defaults().String("pass.time"))
	// fmt.Println("saved:", clickedCreds.pass.time)
	// sss, _ := time.Parse("2020-07-23 09:18:01.882518 +0800 HKT m=+6.644684739", menuet.Defaults().String("pass.time"))
	// fmt.Println("loading:" + sss.String())
	// sst, _ := time.Parse("2020-07-22 15:33:17.865231 +0800 HKT m=+16.132953837", sss.String())
	// fmt.Println("sst:", sst, reflect.TypeOf(sst))
	config.passLength = menuet.Defaults().Integer("passLength")
}

func setDefaults() {
	menuet.Defaults().SetInteger("passLength", config.passLength)
	menuet.Defaults().SetString("uname.value", clickedCreds.uname.value)
	// menuet.Defaults().SetString("uname.time", clickedCreds.uname.time.String())
	menuet.Defaults().SetString("pass.value", clickedCreds.pass.value)
	menuet.Defaults().SetString("pass.time", clickedCreds.pass.time.String())
	// debugNotification("setting: " + clickedCreds.pass.time.String())
	// st := clickedCreds.pass.time.String()
	// fmt.Println(st, reflect.TypeOf(st))
	// debugNotification("converting: " + clickedCreds.pass.time.String())

}

func init() {
	config := settings{}
	currCreds = credentials{}
	clickedCreds = credentials{}

	isDevVersion()
	getDefaults()
	config.Info()
}

func main() {
	setMenuState()

	app := menuet.App()
	app.Children = menuItems
	app.Name = "Password machine"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = Version
	app.AutoUpdate.Repo = "evilcloud/PWgo"
	app.RunApplication()
}

func passwordMachineText() string {
	if !config.devVersion {
		return "Password machine\t" + Version
	}
	return "Password machine\t" + getEmojis(2)
}

// menu items
func menuItems() []item {
	checkDictionaries()

	// getRandomEmoji()
	currCreds.uname.value = generateUsername()
	currCreds.pass.value = generatePassword()
	clipboard.WriteAll(currCreds.pass.value)
	wow := generateWoW()

	spacer := item{}

	return []item{
		item{
			Text:       passwordMachineText(),
			FontWeight: menuet.WeightBlack,
			Clicked: func() {
				debugNotification("clicked Title")
			},
		},
		item{},
		item{Text: "Username"},
		menuDisplayCredential(currCreds.uname.value, "username"),
		item{Text: "Password (" + strconv.Itoa(config.passLength) + " characters)"},
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
				// menuet.Defaults().SetString("uname", details)
			case "password":
				clickedCreds.pass.value = details
				clickedCreds.pass.time = time.Now()
			}
			setDefaults()
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
	setDefaults()
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
				subSubVersion(),
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

func subSubVersion() item {
	return item{
		Text:     Version,
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
				item{Text: "Password " + strconv.Itoa(config.passLength) + " characters)"},
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

// Adds to Version string the time of compile and trues config.devVersion if version is not changed by idflags
func isDevVersion() {
	if Version == devVersionString {
		t := time.Now()
		Version += t.String()
		config.devVersion = true
	}
}

// Messaging

// Pops up error message in the banner and sends to debugNotification error if such presists
func isError(err error) {
	if err != nil {
		debugNotification(err.Error())
		popupMessage("Error", err.Error())
	}
}

// Prints log and pushes banner if development version. Does nothing if -idflagged as production
func debugNotification(text string) {
	if config.devVersion {
		log.Println(text)
		menuet.App().Notification(menuet.Notification{
			Title:   "Debug notification",
			Message: text,
		})
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

// Returns predetermined number of random emojis
func getEmojis(num int) string {
	var emojis string = ""
	for i := 0; i < num; i++ {
		rand.Seed(time.Now().UnixNano())
		emojiNumber := strconv.Itoa((rand.Intn(64)) + 128640)
		emojis += html.UnescapeString("&#" + emojiNumber + ";")
	}
	return emojis
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
			debugNotification("initial dictionary load")
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
				debugNotification("Hello, sailor!")
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
	debugNotification(ret)
	if ret == "a long while ago" {
		ret = ""
	}
	return ret
}
