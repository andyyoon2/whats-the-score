package lib

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/spf13/viper"
)

// Dispatch an API request and return the response body
func Get(path string) ([]byte, error) {
	slog.Info(fmt.Sprintf("Requesting %s\n", path))

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.balldontlie.io%s", path), nil)
	if err != nil {
		return nil, err
	}

	key := viper.GetString("api_key")
	if key == "" {
		return nil, errors.New("API key not found. Run `wts set-api-key`.")
	}

	req.Header.Add("Authorization", key)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	slog.Debug(string(body))
	return body, nil
}
