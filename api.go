package api

import "time"

type InlayId uint64

type Need struct {
	Name   string
	Amount uint
}

type Choice struct {
	Inning string
	To     InlayId
	Needs  []Need
}

type Inlay struct {
	ID        InlayId   `json:"id"`
	Inning    string    `json:"inning"`
	Choices   []Choice  `json:"choices"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InlayDocument struct {
	Inlay Inlay `json:"inlay"`
}

type IDDocument struct {
	ID int `json:"id"`
}
