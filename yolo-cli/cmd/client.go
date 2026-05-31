package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func apiGet(path string, result interface{}) error {
	resp, err := http.Get(serverURL + "/api/v1" + path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

func apiPost(path string, body interface{}, result interface{}) error {
	data, _ := json.Marshal(body)
	resp, err := http.Post(serverURL+"/api/v1"+path, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

func apiPut(path string, body interface{}, result interface{}) error {
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", serverURL+"/api/v1"+path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Yolo-CLI", "true")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

func apiDelete(path string) error {
	req, _ := http.NewRequest("DELETE", serverURL+"/api/v1"+path, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s", body)
	}

	return nil
}

type apiResponse struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
