package generators

import (
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	t "github.com/evilcloud/PWgo/internal/types"

	"github.com/caseymrm/menuet"
)

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

type generator struct {
	nouns        []string
	adjectives   []string
	config       t.Settings
	menu         *menuet.Application
	specialChars []string
	specialNums  []string
}

// NewGenerator is a generator constructor
func NewGenerator(config t.Settings, adjectives, nouns []string) *generator {
	return &generator{
		nouns:        nouns,
		adjectives:   adjectives,
		config:       config,
		//menu:         menuet.Application{},
		specialChars: strings.Split("! @ # $ % & * - + = ?", " "),
		specialNums:  strings.Split("1 2 3 4 5 6 7 8 9", " "),
	}
}

// func (g *generator) EmojisSeeded(num int, seed time.Time) string {
// 	var emojis string = ""
// 	for i := 0; i < num; i++ {
// 		rand.Seed(seed.AddDate(i, i, i).UnixNano())
// 		emojiNumber := strconv.Itoa((rand.Intn(64)) + 128640)
// 		emojis += html.UnescapeString("&#" + emojiNumber + ";")
// 	}
// 	return emojis
// }

func (g *generator) Username() string {
	return strings.Join(g.generateAdjectiveNounPair(), "")
}

func (g *generator) Password() string {
	// FIXME: put the check in the init part of the code
	if g.config.PassLength < passShort {
		g.config.PassLength = passStandard
	}

	var pass []string
	totalDict := append(g.adjectives, g.nouns...)
	for i := 1; i < 5000; i++ {
		pass = append([]string{g.pickRandomWord(totalDict)}, pass...)
		lenPassAlpha := len(strings.Join(pass, "")) + 2
		if lenPassAlpha == g.config.PassLength {
			pass := g.insertRandomNumChar(pass)
			return strings.Join(pass, "")
		}
		if lenPassAlpha > g.config.PassLength {
			pass = nil
		}
	}
	// debugNotification("pass legth failed")
	// popupMessage("Password failure", "Failed to match pass to length!")
	return "Failed to match pass to length"
}

func (g *generator) WoW() string {
	return strings.Join(g.generateAdjectiveNounPair(), " ")
}

// Returns a joined pair of words from generateAdjectiveNounPair() function
// func generateUsername() string {
// 	return strings.Join(generateAdjectiveNounPair(), "")
// }

// Returns a joined with a space pair of word from generateAdjectiveNounPair() function
// func generateWoW() string {
// 	return strings.Join(generateAdjectiveNounPair(), " ")
// }

// Generates a pair of random adjective and noun (only alpha characters) and returns as an array
// Used mostly by Words of Wisdom and Username functions

func (g *generator) generateAdjectiveNounPair() []string {
	var ret []string
	reg, err := regexp.Compile("[^a-zA-Z]+")
	g.isError(err)

	adj := reg.ReplaceAllLiteralString(strings.Title(g.pickRandomWord(g.adjectives)), "")
	nou := reg.ReplaceAllLiteralString(strings.Title(g.pickRandomWord(g.nouns)), "")
	ret = append(ret, []string{adj, nou}...)
	return ret
}

// // Generates a string of random words from adjectives and nouns dictionaries, and returns as a string
// func generatePassword() string {

// }

// Generates random number (except 0) and some special character (from the list). Returns as an array
func (g *generator) insertRandomNumChar(data []string) []string {
	numeral := g.pickRandomWord(g.specialNums)
	char := g.pickRandomWord(g.specialChars)

	if g.config.RandomPlacing {
		data = g.insertIntoPosition(data, numeral)
		data = g.insertIntoPosition(data, char)
	} else {
		data = append(data, []string{numeral}...)
		data = append(data, []string{char}...)
	}
	return data
}

// Inject a string into a random position in array
func (g *generator) insertIntoPosition(data []string, insertion string) []string {
	// I am really sorry for this loop. I have not figured out why slice concatenation doesn't work
	var newData []string
	dataLength := len(data)
	position := g.pickNumberRange(dataLength + 1)
	if position == dataLength {
		newData = append(data, []string{insertion}...)
	} else {
		for i, entry := range data {
			if i == position {
				newData = append(newData, []string{insertion}...)
			}
			newData = append(newData, entry)
		}
	}
	return newData
}

// // Returns predetermined number of random emojis
// func getEmojis(num int) string {
// 	var emojis string = ""
// 	for i := 0; i < num; i++ {
// 		rand.Seed(time.Now().UnixNano())
// 		emojiNumber := strconv.Itoa((rand.Intn(64)) + 128640)
// 		emojis += html.UnescapeString("&#" + emojiNumber + ";")
// 	}
// 	return emojis
// }

// func getEmojisSeeded(num int, seed time.Time) string {

// }

// Picks a random number from provided numbers
func (g *generator) pickNumberRange(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

// Picks a random item from provided string array
func (g *generator) pickRandomWord(data []string) string {
	rand.Seed(time.Now().UnixNano())
	return strings.Trim(strings.Title(data[rand.Intn(len(data))]), "\n")
}

// Pops up error message in the banner and sends to debugNotification error if such presists
func (g *generator) isError(err error) {
	if err != nil {
		g.debugNotification(err.Error())
		g.popupMessage("Error", err.Error())
	}
}

// Prints log and pushes banner if development version. Does nothing if -idflagged as production
func (g *generator) debugNotification(text string) {
	if g.config.DevVersion {
		log.Println(text)
		// menuet.App().Notification(menuet.Notification{
		// 	Title:   "Debug notification",
		// 	Message: text,
		// })
	}
}

// Creates a banner with the predetermined title and message. Dependency on menuet library
func (g *generator) popupMessage(title, message string) {
	g.menu.Notification(menuet.Notification{
		Title:   title,
		Message: message,
	})
}
