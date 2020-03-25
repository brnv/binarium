package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/kovetskiy/godocs"
	"github.com/kovetskiy/lorg"
)

var (
	logger  = lorg.NewLog()
	version = "[manual build]"
)

const usage = `name

Usage:
    name [options]
    name -h | --help

Options:
    --asset <string>   Asset.
    --sum <string>     Transaction price.
    --jwt <string>     Auth token.
    --up               Set option up.
    --down             Set option down.
    --dry-run          Dry run.
    --debug            Enable debug output.
    --trace            Enable trace output.
    -h --help          Show this help.
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

	jwt := optionJWT
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
	logger.Infof("option[asset]: %s", form.Get("option[asset]"))
	logger.Infof("option[currency]: %s", form.Get("option[currency]"))
	logger.Infof("option[expiration][date]: %s", form.Get("option[expiration][date]"))
	logger.Infof("option[expiration][type]: %s", form.Get("option[expiration][type]"))
	logger.Infof("option[kind]: %s", form.Get("option[kind]"))
	logger.Infof("option[source]: %s", form.Get("option[source]"))
	logger.Infof("option[sum]: %s", form.Get("option[sum]"))
}
