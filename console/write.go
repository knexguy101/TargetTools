package console

import (
	"fmt"
	"time"
)

func Write(message string) {
	fmt.Println(fmt.Sprintf("[LOG] [%s] -> %s", time.Now().Format(time.RFC3339), message))
}

func WriteMultipleLines(title string, message []string) {
	fmt.Println(fmt.Sprintf("-----[%s]-----", title))
	for _, v := range message {
		fmt.Println(v)
	}
	fmt.Println("--------------------")
}
