# Client for [Coincheck Exchange API](https://coincheck.jp/documents/exchange/api)

A simple Coincheck Exchange API client.

Example of usage:

```go
package main

import (
	"fmt"

	"github.com/ivanlemeshev/coincheck"
)

const (
	key    = "KEY"
	secret = "SECRET"
)

func main() {
	api := coincheck.New(key, secret)
	result, err := api.GetBalance()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}

	fmt.Printf("Result: %+v\n", result)
}
```

Todo:
- [x] Ticker
- [x] Public trades
- [x] Order book
- [x] New order
- [ ] Unsettled order list
- [ ] Cancel order
- [ ] Order transactions
- [ ] Positions list
- [x] Balance
- [x] Leverage balance
- [ ] Send BTC
- [ ] BTC sends history
- [ ] BTC deposits history
- [ ] Fast bitcoin deposit
- [x] Account information
- [ ] Bank account list
- [ ] Register bank account
- [ ] Remove bank account
- [ ] Withdraw history
- [ ] Create withdraw
- [ ] Cancel withdraw
- [ ] Apply borrowing
- [ ] Borrowing list
- [ ] Reply
- [ ] Transfers to leverage account
- [ ] Transfers from leverage account
