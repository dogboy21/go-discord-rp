package main

import (
	"github.com/dogboy21/go-discord-rp/connection"
	"time"
)

func main() {
	connection.OpenSocket("INSERT APP ID HERE")
	connection.SetActivity("State", "Details", "small_logo", "Small Text.", "bigger_logo", "BIGGER TEXT.")

	for {
		time.Sleep(30 * time.Minute)
	}
}
