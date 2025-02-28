package main

import (
	"bytes"
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
