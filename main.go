package main

import (
	"bufio"
	"fmt"
	"hamkaran_system/bootcamp/final/project/common"
	"hamkaran_system/bootcamp/final/project/service"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		if common.LoginUser.Username != "" {
			fmt.Print(common.LoginUser.Username + ":" + "$ ")
		} else {
			fmt.Print("$ ")
		}
		input, err := reader.ReadString('\n')
		if err == nil {
			args := strings.Fields(input)
			if len(args) == 0 {
				continue
			}
			if args[0] != "history" {
				service.AddHistory(args)
			}
			if len(args) > 2 && (args[len(args)-2] == ">" || args[len(args)-2] == "1>") {
				fileName := args[len(args)-1]
				args = args[:len(args)-2]
				if builtInCommand(args[0]) {

				} else {
					service.OutputRedirectionCommand(args, fileName)
				}
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
			case "adduser":
				service.AddUser(args[1:])
			case "logout":
				service.Logout(args[1:])
			case "history":
				service.HistoryCommand(args[1:])
				//feature 8
			default:
				service.ExecuteCommand(args)
			}
		} else {
			fmt.Println("cannot read input")
		}

	}
}

func builtInCommand(command string) bool {
	switch command {
	case "exit", "echo", "cat", "pwd", "type", "cd", "login", "adduser", "logout", "history":
		return true
	default:
		return false
	}
}
