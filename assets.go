package main

var assets = map[string]string{
	"EUR/USD": "1",
	"USD/JPY": "2",
	"AUD/USD": "3",
	"USD/CAD": "4",
	"EUR/GBP": "7",
	"GBP/USD": "8",
	"NZD/USD": "10",

	"USD/CHF": "11",
	"GBP/JPY": "13",
	"GBP/CHF": "14",
	"AUD/CHF": "15",
	"EUR/CAD": "16",
	"AUD/JPY": "17",
	"CAD/JPY": "18",
	"AUD/NZD": "19",
	"EUR/AUD": "20",

	"GBP/CAD": "21",
	"NZD/JPY": "22",
	"AUD/CAD": "23",
	"EUR/NZD": "24",
	"EUR/CHF": "25",

	"CAD/CHF": "35",
	"CHF/JPY": "36",
	"GBP/AUD": "37",
	"GBP/NZD": "38",
	"NZD/CAD": "39",
	"NZD/CHF": "40",
}

var assetsReverse = map[int]string{
	1:  "EUR/USD",
	2:  "USD/JPY",
	3:  "AUD/USD",
	4:  "USD/CAD",
	7:  "EUR/GBP",
	8:  "GBP/USD",
	10: "NZD/USD",

	11: "USD/CHF",
	13: "GBP/JPY",
	14: "GBP/CHF",
	15: "AUD/CHF",
	16: "EUR/CAD",
	17: "AUD/JPY",
	18: "CAD/JPY",
	19: "AUD/NZD",
	20: "EUR/AUD",

	21: "GBP/CAD",
	22: "NZD/JPY",
	23: "AUD/CAD",
	24: "EUR/NZD",
	25: "EUR/CHF",

	35: "CAD/CHF",
	36: "CHF/JPY",
	37: "GBP/AUD",
	38: "GBP/NZD",
	39: "NZD/CAD",
	40: "NZD/CHF",
}
