package Eldrlang

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println(LOGO)
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>")
		text, _ := input.ReadString('\n')
	}
}
