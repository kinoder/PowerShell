package main

import (
	"bufio"
	"fmt"
	"os"
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
			
		} else {
			fmt.Println("cannot read input")
		}

	}
}

func exitCommand() {

}
