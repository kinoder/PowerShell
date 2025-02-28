package main

import (
	"bytes"
	"fmt"
	"hamkaran_system/bootcamp/final/project/common"
	models "hamkaran_system/bootcamp/final/project/model"
	"hamkaran_system/bootcamp/final/project/service"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestExitCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expectOut  string
		expectExit int
	}{
		{"No argument", []string{}, "exit status 0", 0},
		{"Explicit exit 0", []string{"0"}, "", 0},
		{"Valid exit status", []string{"2"}, "", 2},
		{"Max int exit status", []string{"2147483647"}, "", 2147483647},
		{"Negative exit status", []string{"-1"}, "", -1},
		{"Whitespace in number", []string{" 1 "}, "invalid exit status argument", -1},
		{"Float exit status", []string{"1.5"}, "invalid exit status argument", -1},
		{"Invalid exit status", []string{"abc"}, "invalid exit status argument", -1},
		{"Too many arguments", []string{"1", "2"}, "too many arguments", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestHelperProcess"}, tt.args...)...)
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out
			err := cmd.Run()

			if tt.expectExit >= 0 {
				if exitError, ok := err.(*exec.ExitError); ok {
					if exitError.ExitCode() != tt.expectExit {
						t.Errorf("expected exit code %d, got %d", tt.expectExit, exitError.ExitCode())
					}
				} else if tt.expectExit != 0 {
					t.Errorf("expected exit code %d, but command did not exit", tt.expectExit)
				}
			}

			if tt.expectOut != "" {
				output := strings.TrimSpace(out.String())
				if output != tt.expectOut {
					t.Errorf("expected output %q, got %q", tt.expectOut, output)
				}
			}
		})
	}
}

func TestEchoCommand(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		expectOut string
	}{
		{"Simple text", []string{"Hello", "World"}, "Hello World"},
		{"Text with double quotes", []string{"\"Hello World\""}, "Hello World"},
		{"Text with single quotes", []string{"'Hello World'"}, "Hello World"},
		{"Text with environment variable", []string{"$PATH"}, os.Getenv("PATH")},
		{"Text with escaped character", []string{"\"Hello\\nWorld\""}, "Hello\nWorld"},
		{"Text with multiple spaces", []string{"Hello", "    World"}, "Hello     World"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			service.EchoCommand(tt.arguments)
			w.Close()
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			os.Stdout = old
			output := buf.String()
			if output != tt.expectOut+"\n" {
				t.Errorf("expected output %q, got %q", tt.expectOut, output)
			}
		})
	}
}

func TestCatCommand(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		setupFile func()
		expectOut string
	}{
		{"Empty file", []string{"empty.txt"}, func() {
			os.WriteFile("empty.txt", []byte{}, 0644)
		}, ""},
		{"Single line text", []string{"single.txt"}, func() {
			os.WriteFile("single.txt", []byte("Hello, World!"), 0644)
		}, "Hello, World!"},
		{"Multiple lines", []string{"multiple.txt"}, func() {
			os.WriteFile("multiple.txt", []byte("Hello\nWorld\nThis is a test"), 0644)
		}, "Hello\nWorld\nThis is a test"},
		{"File with special characters", []string{"special.txt"}, func() {
			os.WriteFile("special.txt", []byte("Line 1: $HOME\nLine 2: `echo`"), 0644)
		}, "Line 1: $HOME\nLine 2: `echo`"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFile()
			defer os.Remove(tt.arguments[0])
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			service.CatCommand(tt.arguments)
			w.Close()
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			os.Stdout = old
			output := buf.String()
			if output != tt.expectOut+"\n" {
				t.Errorf("expected output %q, got %q", tt.expectOut, output)
			}
		})
	}
}

func TestTypeCommand(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		setupFile func()
		expectOut string
	}{
		{"Internal command: echo", []string{"echo"}, nil, "echo is a shell builtin\n"},
		{"Internal command: cat", []string{"cat"}, nil, "cat is a shell builtin\n"},
		{"File exists in PATH", []string{"ls"}, nil, "ls: command not found\n"},
		{"Command not found", []string{"nonexistentcommand"}, nil, "nonexistentcommand: command not found\n"},
		{"Missing argument", []string{}, nil, "missing arguments or too many arguments\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFile != nil {
				tt.setupFile()
			}

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			service.TypeCommand(tt.arguments)

			w.Close()
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			os.Stdout = old

			output := buf.String()

			if output != tt.expectOut {
				t.Errorf("expected output %q, got %q", tt.expectOut, output)
			}
		})
	}
}

func TestPwdCommand(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		expectOut string
		expectErr bool
	}{
		{
			"Pwd with no arguments",
			[]string{},
			"c:\\Users\\Armin\\Desktop\\Project\n",
			false,
		},
		{
			"Pwd with arguments",
			[]string{"someArg"},
			"pwd command does not have any argument\n",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldOut := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() { os.Stdout = oldOut }()
			service.PwdCommand(tt.arguments)
			w.Close()
			var buf bytes.Buffer
			buf.ReadFrom(r)
			actualOutput := buf.String()
			if tt.expectErr {
				if actualOutput != tt.expectOut {
					t.Errorf("expected error %q, got %q", tt.expectOut, actualOutput)
				}
			} else {
				if actualOutput != tt.expectOut {
					t.Errorf("expected output %q, got %q", tt.expectOut, actualOutput)
				}
			}
		})
	}
}

func TestCdCommand(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		expectOut string
		expectErr bool
	}{
		{
			"Cd with no arguments",
			[]string{},
			"please enter a path\n",
			true,
		},
		{
			"Cd with too many arguments",
			[]string{"path1", "path2"},
			"too many arguments\n",
			true,
		},
		{
			"Cd to a valid directory",
			[]string{"C:\\Users\\Armin\\Desktop\\Project"},
			"",
			false,
		},
		{
			"Cd to a non-existent directory",
			[]string{"nonexistent_directory"},
			"directory nonexistent_directory does not exist\n",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldOut := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() { os.Stdout = oldOut }()

			service.CdCommand(tt.arguments)

			w.Close()
			var buf bytes.Buffer
			buf.ReadFrom(r)

			actualOutput := buf.String()
			if tt.expectErr {
				if actualOutput != tt.expectOut {
					t.Errorf("expected error %q, got %q", tt.expectOut, actualOutput)
				}
			} else {
				if actualOutput != tt.expectOut {
					t.Errorf("expected output %q, got %q", tt.expectOut, actualOutput)
				}
			}
		})
	}
}

func TestLoginCommand(t *testing.T) {
	users := []models.User{
		{Username: "validUser", Password: "validPassword", ID: 1},
	}

	tests := []struct {
		name             string
		arguments        []string
		expectedUsername string
		expectedErr      bool
	}{
		{
			"Valid Login",
			[]string{"validUser", "validPassword"},
			"validUser",
			false,
		},
		{
			"Invalid Login",
			[]string{"invalidUser", "wrongPassword"},
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			common.LoginUser.Username = ""
			common.LoginUser.ID = 0

			var existingUser models.User
			var err error

			for _, user := range users {
				if user.Username == tt.arguments[0] {
					existingUser = user
					break
				}
			}

			if existingUser.Username == "" || existingUser.Password != tt.arguments[1] {
				err = fmt.Errorf("invalid username or password")
			}

			oldOut := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() { os.Stdout = oldOut }()

			if err != nil {
				fmt.Println("invalid username or password")
			} else {
				common.LoginUser.Username = existingUser.Username
				common.LoginUser.ID = existingUser.ID
			}

			w.Close()
			var buf bytes.Buffer
			buf.ReadFrom(r)

			if tt.expectedErr && common.LoginUser.Username != "" {
				t.Errorf("expected no login but got username: %v", common.LoginUser.Username)
			} else if !tt.expectedErr && common.LoginUser.Username != tt.expectedUsername {
				t.Errorf("expected username %v, but got %v", tt.expectedUsername, common.LoginUser.Username)
			}
		})
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		t.Skip("Skipping helper process test")
	}

	args := os.Args[2:]
	if len(args) > 1 {
		println("too many arguments")
		os.Exit(1)
	}

	if len(args) == 0 {
		println("exit status 0")
		os.Exit(0)
	}

	status, err := strconv.Atoi(args[0])
	if err != nil {
		println("invalid exit status argument")
		os.Exit(1)
	}

	os.Exit(status)
}
