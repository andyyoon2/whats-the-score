package lib

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Dispatch an API request and return the response body
func Get(path string) []byte {
	// fmt.Printf("Would request %s\n", path)
	// return []byte{}

	fmt.Printf("Requesting %s\n", path)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.balldontlie.io/v1%s", path), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", os.Getenv("WTS_API_KEY"))
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
