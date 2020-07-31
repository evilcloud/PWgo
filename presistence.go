package main

import "github.com/caseymrm/menuet"

// Load previous state
func getDefaults() {
	clickedCreds.uname.value = menuet.Defaults().String("uname.value")
	clickedCreds.pass.value = menuet.Defaults().String("pass.value")
	config.passLength = menuet.Defaults().Integer("passLength")
	config.randomPlacing = menuet.Defaults().Boolean("randomPlacing")
}

// Save the state
func setDefaults() {
	menuet.Defaults().SetInteger("passLength", config.passLength)
	menuet.Defaults().SetString("uname.value", clickedCreds.uname.value)
	menuet.Defaults().SetString("pass.value", clickedCreds.pass.value)
	menuet.Defaults().SetBoolean("randomPlacing", config.randomPlacing)
}
