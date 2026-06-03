package main

import (
	"fmt"
	"os"

	snowflake "github.com/Shudharshan07/idflake"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Print("invaid")
		return
	}
	s := &snowflake.SnowFlake{}

	switch args[1] {
	case "init":
		s = snowflake.NewSnowFlake(1)
		fmt.Print("Success")
	case "gen":
		fmt.Print(s.Generate().Base64())
	default:
		break
	}
}
