package main

import (
	"fmt"
	"github.com/shopspring/decimal"
)

func main() {
	d, err := decimal.NewFromString("0.01634790")
	fmt.Printf("Value: %s, Error: %v\n", d.String(), err)
}