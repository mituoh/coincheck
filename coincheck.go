package coincheck

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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

// Balance represents account balance
type Balance struct {
	BTC          string `json:"btc"`
	BTCDebt      string `json:"btc_debt"`
	BTCLendInUse string `json:"btc_lend_in_use"`
	BTCLent      string `json:"btc_lent"`
	BTCReserved  string `json:"btc_reserved"`
	JPY          string `json:"jpy"`
	JPYDebt      string `json:"jpy_debt"`
	JPYLendInUse string `json:"jpy_lend_in_use"`
	JPYLent      string `json:"jpy_lent"`
	JPYReserved  string `json:"jpy_reserved"`
	Success      bool   `json:"success"`
	Error        string `json:"error"`
}

// New creates a new Kraken API struct
func New(key, secret string) (client *APIClient) {
	client = new(APIClient)
	client.key = key
	client.secret = secret
	client.client = new(http.Client)
	return client
}

// ReadBalance returns account balance
func (api APIClient) ReadBalance() (balance Balance, err error) {
	endpoint := URL + "/api/accounts/balance"
	headers := headers(api.key, api.secret, endpoint, "")
	resp, err := api.doRequest("GET", endpoint, headers)
	if err != nil {
		return balance, err
	}
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return balance, err
	}
	if !balance.Success {
		return balance, errors.New(balance.Error)
	}
	return balance, err
}

// doRequest executes a HTTP request to the Coincheck API and returns the result
func (api *APIClient) doRequest(method, endpoint string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, nil)
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
	return hex.EncodeToString(h.Sum(nil))
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
