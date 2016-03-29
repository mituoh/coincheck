package coincheck

import (
	"bytes"
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

// Ticker represents Coincheck ticker
type Ticker struct {
	Ask       float64 `json:"ask"`
	Bid       float64 `json:"bid"`
	High      float64 `json:"high"`
	Last      float64 `json:"last"`
	Low       float64 `json:"low"`
	Timestamp float64 `json:"timestamp"`
	Volume    string  `json:"volume"`
}

// Trade represents Coincheck trade
type Trade struct {
	Amount    string  `json:"amount"`
	CreatedAt string  `json:"created_at"`
	ID        int     `json:"id"`
	OrderType string  `json:"order_type"`
	Rate      float64 `json:"rate"`
}

// OrderBook represents account Coincheck order book
type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
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

// LeverageBalance represents account leverage balance
type LeverageBalance struct {
	Margin struct {
		JPY float64 `json:"jpy"`
	} `json:"margin"`
	MarginAvailable struct {
		JPY string `json:"jpy"`
	} `json:"margin_available"`
	MarginLevel string `json:"margin_level"`
	Success     bool   `json:"success"`
	Error       string `json:"error"`
}

// Accounts represents Coincheck accounts
type Accounts struct {
	BitcoinAddress  string `json:"bitcoin_address"`
	Email           string `json:"email"`
	ID              int    `json:"id"`
	IdentityStatus  string `json:"identity_status"`
	LendingLeverage string `json:"lending_leverage"`
	MakerFee        string `json:"maker_fee"`
	Success         bool   `json:"success"`
	TakerFee        string `json:"taker_fee"`
	Error           string `json:"error"`
}

// Order represents a new order.
type Order struct {
	ID        int     `json:"id"`
	Rate      float64 `json:"rate"`
	Amount    float64 `json:"amount"`
	OrderType string  `json:"order_type"`
	Pair      string  `json:"pair"`
	Success   bool    `json:"success"`
	Error     string  `json:"error"`
}

// New creates a new Kraken API struct
func New(key, secret string) (client *APIClient) {
	client = new(APIClient)
	client.key = key
	client.secret = secret
	client.client = new(http.Client)
	return client
}

// GetTicker returns Coincheck ticker
func (api APIClient) GetTicker() (ticker Ticker, err error) {
	err = api.doGetRequest("/api/ticker", []byte(""), &ticker)
	if err != nil {
		return ticker, err
	}
	return ticker, nil
}

// GetTrades returns Coincheck trades
func (api APIClient) GetTrades() (trades []Trade, err error) {
	err = api.doGetRequest("/api/trades", []byte(""), &trades)
	if err != nil {
		return trades, err
	}
	return trades, nil
}

// GetOrderBook returns Coincheck order book
func (api APIClient) GetOrderBook() (orderBook OrderBook, err error) {
	err = api.doGetRequest("/api/order_books", []byte(""), &orderBook)
	if err != nil {
		return orderBook, err
	}
	return orderBook, nil
}

// GetBalance returns account balance
func (api APIClient) GetBalance() (balance Balance, err error) {
	err = api.doGetRequest("/api/accounts/balance", []byte(""), &balance)
	if err != nil {
		return balance, err
	}
	if !balance.Success {
		return balance, errors.New(balance.Error)
	}
	return balance, nil
}

// GetLeverageBalance returns account leverage balance
func (api APIClient) GetLeverageBalance() (leverageBalance LeverageBalance, err error) {
	err = api.doGetRequest("/api/accounts/leverage_balance", []byte(""), &leverageBalance)
	if err != nil {
		return leverageBalance, err
	}
	if !leverageBalance.Success {
		return leverageBalance, errors.New(leverageBalance.Error)
	}
	return leverageBalance, nil
}

// GetAccounts returns accounts
func (api APIClient) GetAccounts() (accounts Accounts, err error) {
	err = api.doGetRequest("/api/accounts", []byte(""), &accounts)
	if err != nil {
		return accounts, err
	}
	if !accounts.Success {
		return accounts, errors.New(accounts.Error)
	}
	return accounts, nil
}

// NewOrder sends a new order.
func (api APIClient) NewOrder(order Order) (newOrder Order, err error) {
	newOrder = order
	data, err := json.Marshal(newOrder)
	if err != nil {
		return newOrder, err
	}
	err = api.doPostRequest("/api/exchange/orders", data, &newOrder)
	if err != nil {
		return newOrder, err
	}
	if newOrder.Error != "" {
		return newOrder, errors.New(newOrder.Error)
	}
	return newOrder, nil
}

func (api *APIClient) doGetRequest(endpoint string, body []byte, data interface{}) (err error) {
	headers := headers(api.key, api.secret, URL+endpoint, string(body))
	resp, err := api.doRequest("GET", URL+endpoint, body, headers)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, data)
	if err != nil {
		return err
	}
	return nil
}

func (api *APIClient) doPostRequest(endpoint string, body []byte, data interface{}) (err error) {
	headers := headers(api.key, api.secret, URL+endpoint, string(body))
	resp, err := api.doRequest("POST", URL+endpoint, body, headers)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, data)
	if err != nil {
		return err
	}
	return nil
}

func (api *APIClient) doRequest(method, endpoint string, data []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
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
