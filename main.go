package main

import (
	"log"
	"os"

	"github.com/maatko/gowind/tailwind"
)

const (
	TAILWIND_BINARY  = "tailwindcss"
	TAILWIND_VERSION = "latest"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("gowind <update/clean>")
	}

	switch args[0] {
	case "update":
		{
			err := tailwind.Download(TAILWIND_BINARY, TAILWIND_VERSION)
			if err != nil {
				log.Fatal("failed to download TailwindCSS binary", err)
			}
		}
	case "clean":
		{
			path := tailwind.GetBinary(TAILWIND_BINARY)
			if err := os.Remove(path); err != nil {
				log.Fatal("failed to delete TailwindCSS binary", err)
			}
		}
	default:
		{
			log.Fatal("gowind <update/clean>")
		}
	}
}
