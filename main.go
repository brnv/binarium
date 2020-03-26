package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/kovetskiy/godocs"
	"github.com/kovetskiy/lorg"
)

var (
	logger  = lorg.NewLog()
	version = "[manual build]"
)

const usage = `binarium

Usage:
    binarium --login --email <string> --password <string>
    binarium --asset <string> --sum <string> [options]
    binarium --check-id <string> [options]
    binarium -h | --help

Options:
    --login               Login.
    --email <string>      Email.
    --password <string>   Password.
    --asset <string>      Asset.
    --sum <string>        Transaction price.
    --up                  Set option up.
    --down                Set option down.
    --check-id <string>   Check option id.
    --jwt <string>        Auth token.
    --dry-run             Dry run.
    --debug               Enable debug output.
    --trace               Enable trace output.
    -h --help             Show this help.
`

func main() {
	args := godocs.MustParse(usage, version, godocs.UsePager)

	logger.SetIndentLines(true)

	if args["--debug"].(bool) {
		logger.SetLevel(lorg.LevelDebug)
	}

	if args["--trace"].(bool) {
		logger.SetLevel(lorg.LevelTrace)
	}

	if args["--login"].(bool) {
		email := args["--email"].(string)
		password := args["--password"].(string)

		token, err := login(email, password)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info(token)

		os.Exit(0)
	}

	var (
		optionAsset string
		optionJWT   string
		optionSum   string
	)

	if args["--asset"] != nil {
		optionAsset = args["--asset"].(string)
	}

	if args["--jwt"] != nil {
		optionJWT = args["--jwt"].(string)
	}

	if args["--sum"] != nil {
		optionSum = args["--sum"].(string)
	}

	jwt := optionJWT

	if args["--check-id"] != nil {
		transactionID, err := strconv.Atoi(args["--check-id"].(string))
		if err != nil {
			logger.Fatal(err)
		}

		transaction, err := findTransaction(jwt, transactionID)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("%s", transaction.String())

		os.Exit(0)
	}

	optionUp := args["--up"].(bool)
	optionDown := args["--down"].(bool)
	kind := ""

	if optionUp && !optionDown {
		kind = "1"
	}

	if optionDown && !optionUp {
		kind = "2"
	}

	if kind == "" {
		logger.Fatal("no option kind (up or down) specified")
	}

	asset := assets[optionAsset]
	if asset == "" {
		logger.Fatal("no such asset")
	}

	sum := optionSum
	expirationDate := getExpirationDate()

	form := url.Values{}

	form.Add("option[asset]", asset)
	form.Add("option[kind]", kind)
	form.Add("option[sum]", sum)
	form.Add("option[source]", "1")
	form.Add("option[currency]", "1")
	form.Add("option[expiration][date]", expirationDate)
	form.Add("option[expiration][type]", "2")

	printForm(form)

	dryRun := args["--dry-run"].(bool)

	if dryRun {
		logger.Info("nothing to do, dry run mode")
		os.Exit(0)
	}

	err := request(jwt, form)
	if err != nil {
		logger.Error(err)
	}
}

func getExpirationDate() string {
	// notification are coming at
	// 14th minute
	// 29th minute
	// 44th minute
	// 59th minute

	utcTime := time.Now().UTC()
	day := utcTime.Day()
	hour := utcTime.Hour()
	var minute int
	utcMinutes := utcTime.Minute()

	if utcMinutes >= 55 { // 59th minute
		minute = 15
		hour = hour + 1
	} else if utcMinutes >= 40 { // 44th minute
		minute = 0
		hour = hour + 1
	} else if utcMinutes >= 25 { // 29th minute
		minute = 45
	} else if utcMinutes >= 10 { // 14th minute
		minute = 30
	} else if utcMinutes >= 0 { // 59th minute
		minute = 15
	}

	if hour == 24 {
		hour = 0
		day = day + 1
	}

	expirationDate := fmt.Sprintf(
		"%d-%02d-%02dT%02d:%02d:00.000Z",
		utcTime.Year(),
		utcTime.Month(),
		day,
		hour,
		minute,
	)

	return expirationDate
}

func printForm(form url.Values) {
	info := fmt.Sprintf(
		"option[asset]: %s\n"+
			"option[currency]: %s\n"+
			"option[expiration][date]: %s\n"+
			"option[expiration]type]: %s\n"+
			"option[kind]: %s\n"+
			"option[source]: %s\n"+
			"option[sum]: %s",
		form.Get("option[asset]"),
		form.Get("option[currency]"),
		form.Get("option[expiration][date]"),
		form.Get("option[expiration][type]"),
		form.Get("option[kind]"),
		form.Get("option[source]"),
		form.Get("option[sum]"),
	)

	logger.Info(info)
}
