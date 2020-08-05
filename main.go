package main

import (
	"bufio"
	"fmt"
	"github.com/Onelio/Eldrlang/evaluator"
	"github.com/Onelio/Eldrlang/parser"
	"os"
	"strings"
)

func main() {
	fmt.Println(LOGO)
	var (
		input = bufio.NewReader(os.Stdin)
		comp  = parser.NewParser()
		eval  = evaluator.NewEvaluator()
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
		parsed := comp.ParsePackage(code, "main")
		code = ""
		if parsed.Errors.Len() > 0 {
			fmt.Print(parsed.Errors.String())
			continue
		}

		obj := eval.Evaluate(parsed)
		if obj.Errors.Len() > 0 {
			fmt.Print(obj.Errors.String())
			continue
		}
		fmt.Print(obj)
	}
}
