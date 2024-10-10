package fetcher

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

// GenWorkers generate workders
func GenWorkers(num, wait int) chan<- func() {
	tasks := make(chan func())
	for i := 0; i < num; i++ {
		go func() {
			for f := range tasks {
				f()
				time.Sleep(time.Duration(wait) * time.Second)
			}
		}()
	}
	return tasks
}

// FetchURL returns HTTP response body
func FetchURL(url string) (string, error) {
	var errs []error
	httpProxy := viper.GetString("http-proxy")

	resp, body, err := gorequest.New().Proxy(httpProxy).Get(url).Type("text").End()
	if len(errs) > 0 || resp == nil {
		return "", fmt.Errorf("HTTP error. errs: %v, url: %s", err, url)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP error. errs: %v, status code: %d, url: %s", err, resp.StatusCode, url)
	}
	return body, nil
}

type Certificate struct {
	IssuerCAID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
	ResultCount    int    `json:"result_count"`
}

func fetchCertificates(domain string) []Certificate {
	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)

	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return nil
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status:", resp.Status)
		return nil
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	// Unmarshal the JSON response
	var certificates []Certificate
	err = json.Unmarshal(body, &certificates)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)

	}
	return certificates
}
