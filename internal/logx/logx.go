package logx

import "fmt"

const (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

func Watch(msg string) {
	fmt.Println(blue + "[WATCH] " + msg + reset)
}

func Build(msg string) {
	fmt.Println(yellow + "[BUILD] " + msg + reset)
}

func Server(msg string) {
	fmt.Println(green + "[SERVER] " + msg + reset)
}

func Error(msg string) {
	fmt.Println(red + "[ERROR] " + msg + reset)
}