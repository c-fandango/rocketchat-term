package requests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var Host = ""
var Token = ""
var User = ""
var Timeout = 10

var client = &http.Client{
	Timeout: time.Duration(Timeout) * time.Second,
}

func GetRequest(endpoint string, params []map[string]string) ([]byte, error) {

	u := url.URL{Scheme: "https", Host: Host, Path: endpoint}
	q := u.Query()

	for _, paramSet := range params {
		for key, value := range paramSet {
			q.Set(key, value)
		}
	}
	u.RawQuery = q.Encode()
	fmt.Println(u.String())

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to construct request")
	}

	req.Header.Add("X-Auth-Token", Token)
	req.Header.Add("X-User-Id", User)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request returned an error code %v", resp.Status)
	}

	return body, nil
}
