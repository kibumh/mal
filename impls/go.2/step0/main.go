package main

import (
	"fmt"
	"log"

	"github.com/chzyer/readline"
)

func READ(line string) string {
	return line
}

func EVAL(line string) string {
	return line
}

func PRINT(line string) string {
	return line
}

func rep(line string) string {
	return string(line)
}

func main() {
	rl, err := readline.New("user> ")
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for true {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		fmt.Println(rep(string(line)))
	}
}
