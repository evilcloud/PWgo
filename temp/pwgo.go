package main

import (
	"github.com/caseymrm/menuet"
)

func uselessFunction() {}

func menuItems() {
	items := make([]menuet.MenuItem, 0, 42)
	items = append()
}

func main() {
	go uselessFunction()
	app := menuet.App()
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.pwgo"
	app.Children = menuItems
}
