package main

import (
	"fmt"
	"log"

	. "mal/printer"
	. "mal/reader"
	. "mal/types"

	"github.com/chzyer/readline"
)

func READ(line string) (MalType, error) {
	return ReadStr(line)
}

func EVAL(mv MalType) MalType {
	return mv
}

func PRINT(mv MalType) string {
	return PrintStr(mv, true)
}

func rep(line string) (string, error) {
	mv, err := READ(line)
	if err != nil {
		return err.Error(), nil
	}
	return PRINT(EVAL(mv)), nil
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
		mv, _ := rep(string(line))
		fmt.Println(mv)
	}
}
