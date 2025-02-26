package service

import (
	"fmt"
	common "hamkaran_system/bootcamp/final/project/common"
	"hamkaran_system/bootcamp/final/project/database"
	models "hamkaran_system/bootcamp/final/project/model"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var DB = database.DB

// feature 1
func ExitCommand(arguments []string) {
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
func EchoCommand(arguments []string) {
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
func CatCommand(arguments []string) {
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
func TypeCommand(arguments []string) {
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

// feature 5,8
func ExecuteCommand(arguments []string) {
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
func PwdCommand(arguments []string) {
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
func CdCommand(arguments []string) {
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
func LoginCommand(arguments []string) {
	if len(arguments) != 2 {
		fmt.Println("invalid arguments")
		return
	}
	user := &models.User{Username: arguments[0], Password: arguments[1]}
	var existingUser models.User
	err := DB.Where("username = ?", user.Username).First(&existingUser).Error
	if err != nil {
		fmt.Println("invalid username or password")
		return
	}
	err = CheckPassword(arguments[1], existingUser.Password)
	if err != nil {
		fmt.Println("invalid username or password")
		return
	}
	common.LoginUser.Username = user.Username
}

// feature 9
func AddUser(arguments []string) {
	if len(arguments) != 2 {
		fmt.Println("invalid arguments")
		return
	}
	var user = &models.User{
		Username: arguments[0],
		Password: arguments[1],
	}
	err := CreateUser(DB, user)
	if err != nil {
		fmt.Printf("cannot create user : %v\n", err)
		return
	}
	fmt.Println("user created successfully")
}

// feature 9
func Logout(arguments []string) {
	if len(arguments) != 0 {
		fmt.Println("invalid argument")
		return
	}
	common.LoginUser.Username = ""
}

// feature 10
func HistoryCommand(arguments []string) {
	if len(arguments) > 1 {
		fmt.Println("invalid arguments")
		return
	}
	if len(arguments) == 1 {
		if arguments[0] != "clear" {
			fmt.Println("invalid argument")
			return
		}
		if arguments[0] == "clear" {
			ClearHistory()
			return
		}
	}
	ShowHistory()
}

// feature 10
func ClearHistory() {
	//clear history for unlogged users
	if common.LoginUser.Username == "" {
		common.LogHistory = make([]models.LogHistory, 0)
	}
}

// feature 10
func ShowHistory() {
	//show history for unlogged users
	if common.LoginUser.Username == "" {
		if len(common.LogHistory) == 0 {
			fmt.Println("empty command history")
			return
		}
		sort.Slice(common.LogHistory, func(i, j int) bool {
			if common.LogHistory[i].Count == common.LogHistory[j].Count {
				return common.LogHistory[i].CreatedAt.After(common.LogHistory[j].CreatedAt)
			}
			return common.LogHistory[i].Count > common.LogHistory[j].Count
		})
		fmt.Println("| Command | Count |")
		for _, v := range common.LogHistory {
			fmt.Printf("| %s | %d |\n", v.Command, v.Count)
		}
	}
}

// feature 10
func AddHistory(command string) {
	//adding to hsitory for unlogged users
	if common.LoginUser.Username == "" {
		for i, v := range common.LogHistory {
			if v.Command == command {
				common.LogHistory[i].Count++
				return
			}
		}
		common.LogHistory = append(common.LogHistory, models.LogHistory{
			Command:   command,
			Count:     1,
			CreatedAt: time.Now(),
		})
	}
}
