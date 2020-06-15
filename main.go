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

func main() {
	flag.Parse()

	toCurrency := []string{}
	{
		toList := strings.ToUpper(*to)
		for i := 0; i < len(*to)/3; i++ {
			toCurrency = append(toCurrency, string(toList[:3]))
			toList = toList[3:]
		}
	}

	if len(*from) != 3 {
		log.Println("from must be 3 char")
		flag.PrintDefaults()
		os.Exit(1)
	}
	for i := range toCurrency {
		if len(toCurrency[i]) != 3 {
			log.Println("all to currency must be 3 char")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	url := fmt.Sprintf("http://freecurrencyrates.com/api/action.php?do=cvals&f=%s&iso=%s", *from, *to)
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

	for i := range toCurrency {
		fromCurrencyCode, toCurrncyCode, exchangeRate := strings.ToUpper(*from), toCurrency[i], data[toCurrency[i]].(float64)
		if *reverse {
			fromCurrencyCode, toCurrncyCode, exchangeRate = toCurrncyCode, fromCurrencyCode, 1/exchangeRate
		}

		ln, err := json.Marshal(struct {
			FromCurrencyCode string
			ToCurrncyCode    string
			ExchangeRate     float64
			BaseDate         time.Time
		}{fromCurrencyCode, toCurrncyCode, exchangeRate, updatedTime})
		if err != nil {
			log.Println("data create fail", err)
			os.Exit(1)
		}
		fmt.Println(string(ln))
	}
}
