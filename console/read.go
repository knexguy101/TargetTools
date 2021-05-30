package console

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func ReadSingleLine() (string, error) {
	fmt.Println()
	fmt.Printf("[INPUT] > ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	} else {
		return "", errors.New("could not read scanner")
	}
}

func ReadPromptedLine(prompt string) (string, error) {
	fmt.Println()
	fmt.Printf(fmt.Sprintf("[INPUT] [%s] > ", prompt))
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	} else {
		return "", errors.New("could not read scanner")
	}
}
