package main

import (
	"github.com/Dimadetected/chatrooms/internal/server"
	"log"
)

func main() {

	r := server.New()

	if err := r.Run("192.168.1.109:8080"); err != nil {
		log.Fatalln(err)
	}
}
