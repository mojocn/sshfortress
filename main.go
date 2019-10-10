package main

import "sshfortress/cmd"

var buildTime, gitHash string

func main() {
	cmd.Execute(buildTime, gitHash)
}
