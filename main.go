package main

import (
	"bufio"
	"fmt"
	"hamkaran_system/bootcamp/final/project/service"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		//input += " "
		//input2 := "type go.exe"
		//var err error = nil
		if err == nil {
			args := strings.Fields(input)
			if len(args) == 0 {
				continue
			}
			switch args[0] {
			case "exit":
				service.ExitCommand(args[1:])
			case "echo":
				service.EchoCommand(args[1:])
			case "cat":
				service.CatCommand(args[1:])
			case "type":
				service.TypeCommand(args[1:])
			case "pwd":
				service.PwdCommand(args[1:])
			case "cd":
				service.CdCommand(args[1:])
			case "login":
				service.LoginCommand(args[1:])
				//feature 8
			default:
				service.ExecuteCommand(args)
			}
		} else {
			fmt.Println("cannot read input")
		}

	}
}
