package utils

import (
	"io"
	"net/http"
)

func FetchTopLevelPage(url string, c chan []byte) {
	resp, err := http.Get(url)
	if err != nil {
		body := []byte("")
		c <- body
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		body = []byte("")
	}

	c <- body
}
