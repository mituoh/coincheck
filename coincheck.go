package coincheck

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// URL is a Coincheck API base URL
const URL = "https://coincheck.jp"

// APIClient struct represents Coincheck API client
type APIClient struct {
	key    string
	secret string
	client *http.Client
}

// New creates a new Kraken API struct
func New(key, secret string) *APIClient {
	krakenAPI := new(APIClient)
	krakenAPI.key = key
	krakenAPI.secret = secret
	krakenAPI.client = new(http.Client)
	return krakenAPI
}

// ReadBalance
func (api APIClient) ReadBalance() (interface{}, error) {
	endpoint := URL + "api/accounts/balance"
	headers := headers(api.key, api.secret, endpoint, "")
	resp, err := api.doRequest(endpoint, headers)
	return resp, err
}

// doRequest executes a HTTP request to the Coincheck API and returns the result
func (api *APIClient) doRequest(endpoint string, headers map[string]string) (interface{}, error) {
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, requestError(err.Error())
	}
	setHeaders(req, headers)
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, requestError(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, requestError(err.Error())
	}
	return body, nil
}

// headers
func headers(key, secret, uri, body string) map[string]string {
	currentTime := time.Now().Unix()
	nonce := strconv.Itoa(int(currentTime))
	message := nonce + uri + body
	signature := computeHmac256(message, secret)
	headers := map[string]string{
		"Content-Type":     "application/json",
		"ACCESS-KEY":       key,
		"ACCESS-NONCE":     nonce,
		"ACCESS-SIGNATURE": signature,
	}
	return headers
}

// computeHmac256 calculate hash of message usign HMAC SHA256
func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// requestError formats request error
func requestError(err interface{}) error {
	return fmt.Errorf("Could not execute request! (%s)", err)
}

// setHeaders sets request headers
func setHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}
