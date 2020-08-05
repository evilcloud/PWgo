package types

type Settings struct {
	PassLength    int
	RandomPlacing bool
	LoadDict      bool
	DevVersion    bool
	Profanity     struct {
		Sfw    bool
		Nsfw   bool
		Sailor bool
	}
}
