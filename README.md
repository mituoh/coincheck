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
