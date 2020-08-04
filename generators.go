package main

import (
	"html"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
		debugNotification("passLength changed to " + strconv.Itoa(config.passLength))
	}

	var pass []string
	totalDict := append(adjectives, nouns...)
	for i := 1; i < 5000; i++ {
		pass = append([]string{pickRandomWord(totalDict)}, pass...)
		lenPassAlpha := len(strings.Join(pass, "")) + 2
		if lenPassAlpha == config.passLength {
			pass := insertRandomNumChar(pass)
			return strings.Join(pass, "")
		}
		if lenPassAlpha > config.passLength {
			pass = nil
		}
	}
	debugNotification("pass legth failed")
	popupMessage("Password failure", "Failed to match pass to length!")
	return "Failed to match pass to length"
}

// Generates random number (except 0) and some special character (from the list). Returns as an array
func insertRandomNumChar(data []string) []string {
	numeral := pickRandomWord(strings.Split("1 2 3 4 5 6 7 8 9", " "))
	char := pickRandomWord(strings.Split("! @ # $ % & * - + = ?", " "))

	if config.randomPlacing {
		data = insertIntoPosition(data, numeral)
		data = insertIntoPosition(data, char)
	} else {
		data = append(data, []string{numeral}...)
		data = append(data, []string{char}...)
	}
	return data
}

// Inject a string into a random position in array
func insertIntoPosition(data []string, insertion string) []string {
	// I am really sorry for this loop. I have not figured out why slice concatenation doesn't work
	var newData []string
	dataLength := len(data)
	position := pickNumberRange(dataLength + 1)
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

// Picks a random number from provided numbers
func pickNumberRange(num int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(num)
}

// Picks a random item from provided string array
func pickRandomWord(data []string) string {
	rand.Seed(time.Now().UnixNano())
	return strings.Trim(strings.Title(data[rand.Intn(len(data))]), "\n")
}
