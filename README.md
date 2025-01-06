Phada
=====

[![Go Reference](https://pkg.go.dev/badge/github.com/nndi-oss/phada.svg)](https://pkg.go.dev/github.com/nndi-oss/phada)

Phada is a small utility for _dealing with_ AfricasTalking's USSD input. If you've tried to build a USSD application with AT you may have come about the "asterix input problem" i.e. where the data is fed to your application as asterisk separated (ASV?) data.

Phada is a library to reduce the ceremony required to read user's current input, as a bonus you get a way to store the session data

Phada is [Chichewa](https://en.wikipedia.org/wiki/Chichewa) for `Hopscotch`.

## USAGE

Install the module 

```sh
$ go get github.com/nndi-oss/phada@v0.4.0
```

If you're using the standard `net/http` package then you can create the UssdRequestSession
using the `ParseUssdRequest(*http.Request)` function. Otherwise you will have to 
fill the `UssdRequestSession` struct yourself if you're using a framework like Gin, Echo, etc..


```go
import (
    "github.com/nndi-oss/phada"
)

var (
    sessionStore = phada.NewInMemorySessionStore()
)

func handler(w http.ResponseWriter, req *http.Request) {
    session, err := phada.ParseUssdRequest(req)
    if err != nil {
        log.Errorf("Failed to parse request to UssdRequestSession, %s", err)
    }
    err = sessionStore.PutHop(session) // store/persist the session
    if err != nil {
        // handle the error
    }
    session, err = sessionStore.Get(session.SessionID)
    if err != nil {
        log.Errorf("Failed to read UssdRequest from sessionStore, Got Error: %s", err)
    }
    // read the current hop/request text
    currentHopInput := session.ReadIn()
    if currentHopInput == "" {
        fmt.Printf("Failed to read input or input was empty")
    }
    // read the text from the first hop only
    fmt.Printf("Got data: %s during Hop 1", session.GetHopN(1))
    // read text for all the hops (basically the way AT sent it)
    fmt.Printf("Got data: %s", session.Text)
}
```


## LICENSE

MIT License, see [LICENSE.txt](./LICENSE)

---

Copyright (c) 2018 - 2022, NNDI