package main

import (
	"fmt"
	"projet-devops-coveo/cmd"
)

func main() {
	cmdOne := cmd.NewRootCommand()
	err := cmdOne.Execute()
	if err != nil {
		fmt.Println(err)
	}

}
