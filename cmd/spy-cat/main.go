package main

import (
	"log"

	"github.com/illiakornyk/spy-cat/internal/config"
)

func main() {
	config := config.MustLoad()
	log.Printf("Loaded config: %+v\n", config)
}
