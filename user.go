package main

import (
	"chgk/gocket"
	"time"
)

type UserRole int

const (
	UserRoleMaster = UserRole(iota)
	UserRoleCaptain
	UserRoleСonnoisseur
	UserRoleSpectator
)

func NewUser(id, name string, role UserRole, socket *gocket.Socket, online bool) *User {
	return &User{
		ID:             id,
		Name:           name,
		Role:           role,
		Socket:         socket,
		Online:         online,
		ConnectionTime: time.Now().Unix(),
	}
}

type User struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Role           UserRole       `json:"role"`
	Socket         *gocket.Socket `json:"-"`
	Online         bool           `json:"online"`
	ConnectionTime int64          `json:"connection_time"`
}

func (u *User) State() *gocket.EmitterData {
	return &gocket.EmitterData{
		"id":              u.ID,
		"name":            u.Name,
		"role":            u.Role,
		"online":          u.Online,
		"connection_time": u.ConnectionTime,
	}
}
