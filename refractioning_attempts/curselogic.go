package main

type swearState struct {
	SFW bool
	NSFW bool
	Sailor bool
}

signal := swearState{true, false, false}

func aaa(signal string) string {
	switch signal {
	case "NSFW":
		signal.NSFW = true
		signal.SFW = false
		
	}
}