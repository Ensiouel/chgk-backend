package main

import "chgk/gocket"

type UserRole int

const (
	UserRoleMaster = UserRole(iota)
	UserRoleSpectator
	UserRoleCaptain
	UserRoleСonnoisseur
)

func NewUser(id, name string, role UserRole, socket *gocket.Socket, online bool) *User {
	return &User{
		ID:     id,
		Name:   name,
		Role:   role,
		Socket: socket,
		Online: online,
	}
}

type User struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Role   UserRole       `json:"role"`
	Socket *gocket.Socket `json:"-"`
	Online bool           `json:"online"`
}

func (u *User) State() *gocket.EmitterData {
	return &gocket.EmitterData{
		"id":     u.ID,
		"name":   u.Name,
		"role":   u.Role,
		"online": u.Online,
	}
}
