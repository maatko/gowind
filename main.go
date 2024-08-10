package main

import (
	"os"

	"github.com/maatko/gowind/tailwind"
)

func main() {
	tailwind.Execute(os.Args[1:]...)
}
