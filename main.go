package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err == nil {
			args := strings.Fields(input)
			if len(args) == 0 {
				continue
			}
			switch args[0] {
			case "exit":
				{
					exitCommand(args[1:])
				}
			case "echo":
				{
					echoCommand(args[1:])
				}
			default:
				fmt.Println("this command is not recognized as an internal or external command,operable program or batch file.")
			}
		} else {
			fmt.Println("cannot read input")
		}

	}
}

func exitCommand(arguments []string) {
	var status int
	var err error
	if len(arguments) == 0 {
		fmt.Println("exit status 0")
		os.Exit(0)
	}
	if len(arguments) == 1 {
		status, err = strconv.Atoi(arguments[0])
		if err != nil {
			fmt.Println("invalid exit status arguemnt")
			return
		}
		if status == 0 {
			fmt.Println("exit status 0")
			os.Exit(0)
			return
		}
		os.Exit(status)
	} else {
		fmt.Println("too many arguments")
	}
}

func echoCommand(arguments []string) {
	if len(arguments) == 0 {
		fmt.Println()
		return
	}
	result := make([]string, 0, len(arguments))
	for _, arg := range arguments {
		if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") || strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
			result = append(result, arg[1:len(arg)-1])
		} else {
			replaced := os.ExpandEnv(arg)
			result = append(result, replaced)
		}
	}
	fmt.Println(strings.Join(result, " "))
}
