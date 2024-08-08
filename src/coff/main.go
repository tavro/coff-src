package main

import (
	"fmt"
	"os"
	"os/user"
	"coff-src/src/coff/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! I'm the CoffLang interpreter. :-)\n", user.Username)
	fmt.Printf("I'm ready to take your commands!\n")
	repl.Start(os.Stdin, os.Stdout)
}


