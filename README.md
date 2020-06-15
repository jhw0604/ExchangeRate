# ExchangeRate
getting currency exchange rate to jsonl

data based on https://freecurrencyrates.com/


## how to use
from currency to currencys at recent times
ex) KRW to USD and JPY

```
ExchangeRate -from KRW -to USDJPY
```

> {"FromCurrencyCode":"KRW","ToCurrncyCode":"USD","ExchangeRate":0.0008291480646093758,"BaseDate":"2020-06-15T06:00:20Z"}
> {"FromCurrencyCode":"KRW","ToCurrncyCode":"JPY","ExchangeRate":0.08901716260829276,"BaseDate":"2020-06-15T06:00:20Z"}


if need USD to KRW and JPY to KRW then

```
ExchangeRate -from KRW -to USDJPY -reverse
```

you can get

> {"FromCurrencyCode":"USD","ToCurrncyCode":"KRW","ExchangeRate":1206.057208215417,"BaseDate":"2020-06-15T06:00:20Z"}
> {"FromCurrencyCode":"JPY","ToCurrncyCode":"KRW","ExchangeRate":11.233788751505779,"BaseDate":"2020-06-15T06:00:20Z"}

