package pwui

import (
	"github.com/caseymrm/menuet"
)

type item = menuet.MenuItem

func main() {
	app := menuet.App()
	app.Name = "PWgo"
	app.Label = "com.github.evilcloud.PWgo"
	app.AutoUpdate.Version = "v0.2"
	app.AutoUpdate.Repo = "evilcloud.PWgo"
	app.RunApplication()
}

func setMenuState(swearState string) {
	switch swearState{
	case "sfw":
		image := "pw.pdf"
	case "nsfw":
		image := "nsfw.pdf"
		case "sailor" {
			image := "sailor.pdf"
		}

		menuet.App().SetMenuState(&menuet.MenuState{
			Image: image,
		})
	}
	menuet.App().MenuChanged()
}

func menuItems() []item {
	return []item{
		item{Text: "ok"}
	}
}