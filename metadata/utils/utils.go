package utils

import (
	"fmt"
	"net/http"
)

func GetRequest(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении http.Get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDb API вернул ошибку: %s", resp.Status)
	}

	return resp, nil
}
