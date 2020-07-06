package main

import (
	"Eldrlang/parser"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(LOGO)
	var (
		input = bufio.NewReader(os.Stdin)
		p     = parser.NewParser()
		code  = ""
	)
	for {
		fmt.Print(">>")
		line, _ := input.ReadString('\n')
		// Special commands check
		if strings.HasPrefix(line, "exit") {
			return
		}
		if strings.HasPrefix(line, "clear") {
			cleanConsole()
			continue
		}
		// Continue execution
		code += line
		if !strings.Contains(line, ";") {
			continue
		}
		parsed := p.ParsePackage(code, "main")
		code = ""
		if parsed.Errors.Len() > 0 {
			fmt.Print(parsed.Errors.String())
			continue
		}
		fmt.Println(parsed)
	}
}
