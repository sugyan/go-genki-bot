package main

import (
	"flag"
)

func main() {
	user := flag.String("user", "default", "user screen name")
	flag.Parse()
	NewGenkiBot().Run(user)
}
