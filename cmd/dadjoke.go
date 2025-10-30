/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"github.com/spf13/cobra"
)

// dadjokeCmd represents the dadjoke command
var dadjokeCmd = &cobra.Command{
	Use:   "dadjoke",
	Short: "Generate Random Dad Joke",
	Long: `Type dadjoke to generate a random dad joke`,
	Run: func(cmd *cobra.Command, args []string) {
		getRandomDadJoke()
	},
}

func init() {
	rootCmd.AddCommand(dadjokeCmd)

}

type DadJoke struct {
	ID   string `json:"id"`
	Joke string `json:"joke"`
	Status int   `json:"status"`
}

func getRandomDadJoke() {
	baseAPI := "https://icanhazdadjoke.com/"
	responseBytes := getDadJokeData(baseAPI)

	joke := DadJoke{}
	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		fmt.Printf("Failed in parsing data. %v", err)
	}
	fmt.Println(joke.Joke)
}

func getDadJokeData(baseAPI string) []byte {
	req, err := http.NewRequest(http.MethodGet, baseAPI, nil)

	if err != nil {
		fmt.Printf("unlucky! LOL! skill issue: #FAILED TO REQUEST. %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "DadJokeCLI (https://github.com/bradleyalexa/dadjokecli)")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("unlucky! LOL! skill issue: #FAILED TO GET RESPONSE. %v", err)
	}

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read responses. %v", err)
	}

	return responseBytes
}