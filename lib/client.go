package lib

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

// Dispatch an API request and return the response body
func Get(path string) []byte {
	fmt.Printf("Requesting %s\n", path)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.balldontlie.io%s", path), nil)
	if err != nil {
		log.Fatal(err)
	}

	key := viper.GetString("api_key")
	if key == "" {
		slog.Error("API key not found. Run `wts set-api-key`.")
		os.Exit(1)
	}

	req.Header.Add("Authorization", key)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	return body
}
