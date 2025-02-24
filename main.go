package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
				exitCommand(args[1:])
			case "echo":
				echoCommand(args[1:])
			case "cat":
				catCommand(args[1:])
			case "type":
				typeCommand(args[1:])
			case "pwd":
				pwdCommand(args[1:])
			case "cd":
				cdCommand(args[1:])
			case "login":
				loginCommand(args[1:])
				//feature 8
			default:
				executeCommand(args)
			}
		} else {
			fmt.Println("cannot read input")
		}

	}
}

// feature 1
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

// feature 2
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

// feature 3
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

// feature 4
func typeCommand(arguments []string) {
	if len(arguments) != 1 {
		fmt.Println("missing arguments or too many arguments")
		return
	}
	command := arguments[0] + ".exe"
	builtInTypes := map[string]bool{
		"exit": true, "echo": true, "cat": true, "pwd": true, "type": true, "cd": true, "login": true, "adduser": true, "logout": true,
		"history": true,
	}

	for k := range builtInTypes {
		if arguments[0] == k {
			fmt.Printf("%s is a shell builtin\n", arguments[0])
			return
		}
	}
	env := os.Getenv("PATH")
	paths := strings.Split(env, string(os.PathListSeparator))

	for _, dir := range paths {

		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		path := filepath.Join(dir, command)

		info, err := os.Stat(path)
		if err == nil {
			if info.Mode().IsRegular() {
				fmt.Printf("%s is %s\n", arguments[0], path)
				return
			}
		}
	}
	fmt.Printf("%s: command not found\n", arguments[0])
}

// feature 5

func executeCommand(arguments []string) {
	command := arguments[0]
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	var cmdPath string
	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)

		info, err := os.Stat(fullPath + ".exe")
		if err == nil && !info.IsDir() {
			cmdPath = fullPath
			break
		}
	}

	if cmdPath == "" {
		fmt.Printf("%s: command not found\n", command)
		return
	}

	cmd := exec.Command(cmdPath, arguments[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s: %v\n", strings.Join(arguments, " "), err)
	}
}

// feature 6
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

// feature 7
func cdCommand(arguments []string) {
	if len(arguments) == 0 {
		fmt.Println("please enter a path")
		return
	}
	if len(arguments) > 1 {
		fmt.Println("too many arguments")
		return
	}
	dir := arguments[0]
	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("directory %s does not exist\n", dir)
		return
	}
}

// feature 9
func loginCommand(arguments []string) {
	if len(arguments) == 0 {
		fmt.Println("missing arguments")
		return
	}
}
