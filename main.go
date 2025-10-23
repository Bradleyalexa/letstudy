/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package main

import ("github.com/bradleyalexa/letstudy/cmd"
"github.com/bradleyalexa/letstudy/data")
func main() {
	data.OpenDB()
	cmd.Execute()
}
