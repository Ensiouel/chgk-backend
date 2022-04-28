package main

import (
	"chgk/gocket"
	"time"
)

func NewTable() *Table {
	return &Table{
		Users:        map[*User]bool{},
		ScoreBoard:   map[string]int{"Игроки": 1, "Зрители": 3},
		PackID:       "",
		TourID:       0,
		Open:         true,
		Duration:     0,
		TimerRunning: false,
	}
}

type Table struct {
	Users            map[*User]bool
	ScoreBoard       map[string]int
	PackID           string
	TourID           int
	Open             bool
	Duration         int
	TimerRunning     bool
	timerStart       time.Time
	secondsRemaining int64
}

func (t *Table) ContainsUser(id string) bool {
	for user := range t.Users {
		if user.ID == id {
			return true
		}
	}
	return false
}

func (t *Table) GetUser(id string) *User {
	for user := range t.Users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

func (t *Table) State() *gocket.EmitterData {
	users := make([]*User, 0, len(t.Users))

	for user := range t.Users {
		users = append(users, user)
	}

	if t.TimerRunning {
		t.secondsRemaining = int64(t.Duration*1000) - time.Now().Sub(t.timerStart).Milliseconds()
	} else {
		t.secondsRemaining = 0.0
	}

	return &gocket.EmitterData{
		"users":             users,
		"score_board":       t.ScoreBoard,
		"pack_id":           t.PackID,
		"tour_id":           t.TourID,
		"open":              t.Open,
		"duration":          t.Duration,
		"timer_running":     t.TimerRunning,
		"seconds_remaining": t.secondsRemaining,
	}
}
