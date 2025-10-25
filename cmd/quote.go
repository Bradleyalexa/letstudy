/*
Copyright © 2025 NAME HERE
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// quoteCmd represents the quote command
var quoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "Generate Random Quote",
	Long:  "Generate a random quote from zenquotes.io",
	Run: func(cmd *cobra.Command, args []string) {
		getRandomQuote()
	},
}

func init() {
	rootCmd.AddCommand(quoteCmd)
}

type ZenQuote struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
}

func getQuoteData(baseAPI string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, baseAPI, nil)
	if err != nil {
		fmt.Printf("Failed to make request: %v", err)
		//return nil
	}

	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to get response: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", res.Status)
	}

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read response: %v", err)
	}

	return responseBytes, nil
}

func getRandomQuote() {
	baseAPI := "https://zenquotes.io/api/random"
	responseBytes, err := getQuoteData(baseAPI)
	if err != nil {
		fmt.Println(err)
		return
	}

	var quotes []ZenQuote
	if err := json.Unmarshal(responseBytes, &quotes); err != nil {
		fmt.Printf("Failed to parse response: %v\n", err)
		return
	}

	if len(quotes) == 0 {
		fmt.Println("Internal server error - Quote not provided.")
		return
	}

	q := quotes[0]
	fmt.Printf("\"%s\"\n— %s\n", q.Quote, q.Author)
}
