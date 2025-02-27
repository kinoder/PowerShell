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
	joinedArgs := strings.Join(arguments, " ")
	fmt.Println(processArgument(joinedArgs))
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
	common.LoginUser.Username = existingUser.Username
	common.LoginUser.ID = existingUser.ID
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
	if common.LoginUser.Username == "" {
		common.LogHistory = make([]models.LogHistory, 0)
	} else {
		err := DB.Where("user_id = ?", common.LoginUser.ID).Delete(&models.LogHistory{}).Error
		if err != nil {
			fmt.Println("error clearing history:", err)
			return
		}
		fmt.Println("command history cleared for user:", common.LoginUser.Username)
	}
}

// feature 10
func ShowHistory() {
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
			fmt.Printf("| %-15s | %-5d |\n", v.Command, v.Count)
		}
	} else {
		var userHistory []models.LogHistory
		err := DB.Where("user_id = ?", common.LoginUser.ID).Order("count DESC, created_at DESC").Find(&userHistory).Error
		if err != nil {
			fmt.Println("error fetching history:", err)
			return
		}
		if len(userHistory) == 0 {
			fmt.Println("empty command history")
			return
		}
		fmt.Println("| Command         | Count |")
		for _, v := range userHistory {
			fmt.Printf("| %-15s | %-5d |\n", v.Command, v.Count)
		}
	}
}

// feature 10
func AddHistory(command []string) {
	commands := strings.Join(command, " ")
	if common.LoginUser.Username == "" {
		for i, v := range common.LogHistory {
			if v.Command == commands {
				common.LogHistory[i].Count++
				return
			}
		}
		common.LogHistory = append(common.LogHistory, models.LogHistory{
			Command:   commands,
			Count:     1,
			CreatedAt: time.Now(),
		})
	} else {
		var history models.LogHistory
		err := DB.Where("user_id = ? AND command = ?", common.LoginUser.ID, command).First(&history).Error
		if err == nil {
			history.Count++
			history.CreatedAt = time.Now()
			DB.Save(&history)
		} else {
			newHistory := models.LogHistory{
				UserId:    common.LoginUser.ID,
				Command:   commands,
				Count:     1,
				CreatedAt: time.Now(),
			}
			DB.Create(&newHistory)
		}
	}
}

func OutputRedirectionCommand(argumetns []string, fileName string) {
	cmd := exec.Command(argumetns[0], argumetns[1:]...)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	cmd.Stdout = file
	cmd.Stderr = file

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}

func OutputRedirectionBuiltInCommandRed(arguments []string, fileName string) {
	switch arguments[0] {
	case "echo":
		result := strings.Join(arguments[1:], " ")
		writeToFile(result, fileName)
	case "cat":
		if len(arguments) > 1 {

			result := readFile(arguments[1])
			writeToFile(result, arguments[1])
		} else {
			writeToFile("file does not exist", fileName)
		}
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			writeToFile("error getting current directory", fileName)
		} else {
			writeToFile(dir, fileName)
		}
	case "type":
		if len(arguments) > 1 {
			writeToFile(arguments[1], fileName)
		} else {
			writeToFile("missing argument", fileName)
		}
	default:
		writeToFile("unknown command: "+arguments[0], fileName)
	}
}

// //////////////////////////////////////////////
func processArgument(arg string) string {
	if len(arg) > 1 && arg[0] == '"' && arg[len(arg)-1] == '"' {
		return parseDoubleQuoted(arg[1 : len(arg)-1])
	} else if len(arg) > 1 && arg[0] == '\'' && arg[len(arg)-1] == '\'' {
		return arg[1 : len(arg)-1]
	}
	return expandEnvVariables(arg)
}

func parseDoubleQuoted(input string) string {
	var result strings.Builder
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' && i+1 < len(input) {
			switch input[i+1] {
			case '$', '`', '"', '\\':
				i++
				result.WriteByte(input[i])
			case 'n':
				i++
				result.WriteByte('\n')
			default:
				result.WriteByte(input[i])
			}
		} else {
			result.WriteByte(input[i])
		}
	}
	return result.String()
}

func expandEnvVariables(input string) string {
	var result strings.Builder
	var i int
	for i < len(input) {
		if input[i] == '$' && i+1 < len(input) && (isAlpha(input[i+1]) || input[i+1] == '_') {
			j := i + 1
			for j < len(input) && (isAlphaNum(input[j]) || input[j] == '_') {
				j++
			}
			varName := input[i+1 : j]
			value := os.Getenv(varName)
			if value != "" {
				result.WriteString(value)
			} else {
				result.WriteString("$")
				result.WriteString(varName)
			}
			i = j
		} else {
			if input[i] == '\\' && i+1 < len(input) && (input[i+1] == '$' || input[i+1] == '`' || input[i+1] == '"' || input[i+1] == '\\') {
				i++
			}
			result.WriteByte(input[i])
			i++
		}
	}
	return result.String()
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isAlphaNum(c byte) bool {
	return isAlpha(c) || (c >= '0' && c <= '9')
}

func writeToFile(input string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("cannot create file: ", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(input)
	if err != nil {
		fmt.Println("cannot write to file : ", err)
	}
}

func readFile(fileName string) string {
	content, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Println("cannot read the file : ", err)
	}
	return string(content)
}
