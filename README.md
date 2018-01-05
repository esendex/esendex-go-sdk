# esendex

[![Build Status](https://travis-ci.org/esendex/esendex-go-sdk.svg)](https://travis-ci.org/esendex/esendex-go-sdk) [![GoDoc](https://godoc.org/github.com/esendex/esendex-go-sdk?status.svg)](https://godoc.org/github.com/esendex/esendex-go-sdk)

A client for the [Esendex REST API][esendex].

``` bash
$ go get github.com/esendex/esendex-go-sdk
```

Extremely simple example:

``` go
import (
    "github.com/esendex/esendex-go-sdk"
)

var (
    accountReference = "EX000000"
    username         = "user"
    password         = "pass"
)

func main() {
    client := esendex.New(username, password)

    accountClient := client.Account(accountReference)

    response, err := accountClient.Received()
    if err == nil {
        // response.Messages contains messages received by the account
    }
}

```

See the GoDocs for some further examples around sending messages, checking message status etc.

[esendex]: http://developers.esendex.com/APIs/REST-API
