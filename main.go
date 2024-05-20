package main

import (
	"projet-devops-coveo/cmd"

	"github.com/sirupsen/logrus"
)

func main() {
	cmdOne := cmd.NewRootCommand()
	err := cmdOne.Execute()
	if err != nil {
		logrus.Error(err)
	}

}
