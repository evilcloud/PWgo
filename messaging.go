package main

import (
	"github.com/caseymrm/menuet"
)

// // Pops up error message in the banner and sends to debugNotification error if such presists
// func isError(err error) {
// 	if err != nil {
// 		debugNotification(err.Error())
// 		popupMessage("Error", err.Error())
// 	}
// }

// // Prints log and pushes banner if development version. Does nothing if -idflagged as production
// func debugNotification(text string) {
// 	if config.devVersion {
// 		log.Println(text)
// 		// menuet.App().Notification(menuet.Notification{
// 		// 	Title:   "Debug notification",
// 		// 	Message: text,
// 		// })
// 	}
// }

// Creates a banner with the predetermined title and message. Dependency on menuet library
func popupMessage(title, message string) {
	menuet.App().Notification(menuet.Notification{
		Title:   title,
		Message: message,
	})
}
