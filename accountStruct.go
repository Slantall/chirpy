package main

type LoginAccount struct {
	Password           string `json:"password"`
	Email              string `json:"email"`
	Expires_in_seconds int    `json:"expires_in_seconds"`
}
