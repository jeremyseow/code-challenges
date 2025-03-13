package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendWebHook(url, id string, data interface{}) error {
	_, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(content))
	return nil
}
