package main

import (
	"fmt"
)

const DebuggerLevel = 2

func Debug(level int, format string, a ...interface{}) {
	if level <= DebuggerLevel {
		fmt.Print("[DEBUG] ")
		fmt.Printf(format, a...)
		fmt.Println()
	}
}
