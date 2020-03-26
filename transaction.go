package main

import (
	"fmt"
)

type Transaction struct {
	ID         int     `json:"id"`
	Asset      int     `json:"asset"`
	Kind       int     `json:"kind"` // 1 - up, 2 - down
	Sum        float32 `json:"sum"`
	Status     int     `json:"status"` // 1 - in progress, 2 - done
	Income     float32 `json:"income"`
	QuoteOpen  float32 `json:"quoteOpen"`
	QuoteClose float32 `json:"quoteClose"`
}

func (transaction *Transaction) String() string {
	status := ""

	success := ""

	if transaction.Status == 1 {
		status = "in progress"
	} else if transaction.Status == 2 {
		status = "done"
		if transaction.Income != 0 {
			success = "true"
		} else {
			success = "fail"
		}
	}

	profit := transaction.Income - transaction.Sum
	asset := assetsReverse[transaction.Asset]

	result := fmt.Sprintf(
		"ID: %d\n"+
			"Asset: %s\n"+
			"Status: %s\n"+
			"Success: %s\n"+
			"Income: %f\n"+
			"Profit: %f",
		transaction.ID, asset, status, success,
		transaction.Income, profit,
	)

	return result
}
