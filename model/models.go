package models

type User struct {
	Usermail string `json:"username"`
	Password string `json:"password"`
}

var UserStore = map[string]string{}

