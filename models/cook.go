package models

const (
	STOVE = "stove"
	OVEN  = "oven"
)

type Cooks struct {
	Cooks []Cook `json:"cooks"`
}

type Cook struct {
	Name        string `json:"name"`
	Rank        int    `json:"rank"`
	Proficiency int    `json:"proficiency"`
	CatchPhrase string `json:"catch-phrase"`
}

type Apparat struct {
}

type CookingApparatusQueue struct {
	Ovens  chan Apparat
	Stoves chan Apparat
}

type CookingApparatusTypes struct {
	Ovens  string
	Stoves string
}

