package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	from    = flag.String("from", "", "base currncy code")
	to      = flag.String("to", "", "currncy code list ex. JPY, EUR => JPYEUR")
	reverse = flag.Bool("reverse", false, "`to` to `from`")
)

var aliasCode = map[string]string{
	"NTD": "TWD",
}

func main() {
	flag.Parse()

	toCurrency := map[string]struct{}{}
	{
		toList := strings.ToUpper(*to)
		for i := 0; i < len(*to)/3; i++ {
			to := toList[:3]

			alias, ok := aliasCode[to]
			if ok {
				toCurrency[alias] = struct{}{}
			} else {
				toCurrency[to] = struct{}{}
			}
			toList = toList[3:]
		}
	}

	if len(*from) != 3 {
		log.Println("from must be 3 char")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var toCurrencyCode string
	for to := range toCurrency {
		if len(to) != 3 {
			log.Println("all to currency must be 3 char")
			flag.PrintDefaults()
			os.Exit(1)
		}
		toCurrencyCode += to
	}

	url := fmt.Sprintf("http://freecurrencyrates.com/api/action.php?do=cvals&f=%s&iso=%s", *from, toCurrencyCode)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("api call fail", err)
		os.Exit(1)
	}

	buff := bytes.Buffer{}
	buff.ReadFrom(resp.Body)
	data := map[string]interface{}{}

	err = json.Unmarshal(buff.Bytes(), &data)
	if err != nil {
		log.Println("response json parsing fail", err)
		os.Exit(1)
	}
	updated, ok := data["updated"].(string)
	if !ok {
		log.Println("response updated parsing fail", err)
		os.Exit(1)
	}

	unixUpdated, err := strconv.Atoi(updated)
	if err != nil {
		log.Println("response updated convert string to integer fail", err)
		os.Exit(1)
	}

	updatedTime := time.Unix(int64(unixUpdated), 0).UTC()

	reverseAliasCode := map[string]string{}
	for key, value := range aliasCode {
		reverseAliasCode[value] = key
	}

	for toCurrencyCode := range toCurrency {
		fromCurrencyCode, exchangeRate := strings.ToUpper(*from), data[toCurrencyCode].(float64)
		reverseToCurrencyCode, ok := reverseAliasCode[toCurrencyCode]
		if ok {
			printJsonl(fromCurrencyCode, reverseToCurrencyCode, exchangeRate, updatedTime, *reverse)
		}
		printJsonl(fromCurrencyCode, toCurrencyCode, exchangeRate, updatedTime, *reverse)
	}
}

func printJsonl(fromCurrencyCode, toCurrencyCode string, exchangeRate float64, updatedTime time.Time, reverse bool) {
	if reverse {
		fromCurrencyCode, toCurrencyCode, exchangeRate = toCurrencyCode, fromCurrencyCode, 1/exchangeRate
	}
	jsonl, err := json.Marshal(struct {
		FromCurrencyCode string
		ToCurrencyCode   string
		ExchangeRate     float64
		BaseDate         time.Time
	}{fromCurrencyCode, toCurrencyCode, exchangeRate, updatedTime})
	if err != nil {
		log.Println("data create fail", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonl))
}
