package main

import (
	"log"

	"github.com/oswines/landset/internal/hoard"
)

func main() {
	log.Println("Starting server on port 3469")
	hoard.NewHTTPServer(":3469")
}
