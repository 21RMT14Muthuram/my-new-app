package models

type User struct {
	Usermail string `json:"usermail"`
	Password string `json:"password"`
}

var UserStore = map[string]string{}

