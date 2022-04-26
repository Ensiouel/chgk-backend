package main

import "chgk/gocket"

func NewTable() *Table {
	return &Table{
		Users:      []*User{},
		ScoreBoard: map[string]int{"Игроки": 1, "Зрители": 3},
		PackID:     "",
		TourID:     0,
		Open:       true,
	}
}

type Table struct {
	Users      []*User
	ScoreBoard map[string]int
	PackID     string
	TourID     int
	Open       bool
}

func (t *Table) ContainsUser(id string) bool {
	for _, user := range t.Users {
		if user.ID == id {
			return true
		}
	}
	return false
}

func (t *Table) GetUser(id string) *User {
	for _, user := range t.Users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

func (t *Table) State() *gocket.EmitterData {
	return &gocket.EmitterData{
		"users":       t.Users,
		"score_board": t.ScoreBoard,
		"pack_id":     t.PackID,
		"tour_id":     t.TourID,
		"open":        t.Open,
	}
}
