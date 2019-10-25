package main

import "github.com/mojocn/sshfortress/cmd"

var buildTime, gitHash string

func main() {
	cmd.Execute(buildTime, gitHash)
}
