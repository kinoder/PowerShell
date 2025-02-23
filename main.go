package main

import (
	"bufio"
	"fmt"
	"io"
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
			case "cat":
				{
					catCommand(args[1:])
				}
			case "pwd":
				{
					pwdCommand(args[1:])
				}
			default:
				fmt.Println("this command is not recognized as an internal or external command,operable program or batch file.")
			}
		} else {
			fmt.Println("cannot read input")
		}

	}
}
//1
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
//2
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
//3
func catCommand(arguments []string) {
	if len(arguments) == 0 {
		fmt.Println("file does not exist")
		return
	}
	for _, fname := range arguments {
		file, err := os.Open(fname)
		if err != nil {
			fmt.Printf("cannot open %s. error : %v\n", fname, err)
			return
		}
		defer file.Close()
		_, err = io.Copy(os.Stdout, file)
		if err != nil {
			fmt.Printf("cannot read %s. error : %v\n", fname, err)
		}
		fmt.Println()
	}
}
//6
func pwdCommand(arguments []string) {
	if len(arguments) > 0 {
		fmt.Println("pwd command does not have any argument")
		return
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getting current directory")
		return
	}
	fmt.Println(dir)
}
