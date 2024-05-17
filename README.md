# vrc_osc_query

# 
[VRC OSC Query ](https://github.com/vrchat-community/vrc-oscquery-lib) library for Golang.
We use the modified [go-osc](https://github.com/StarLight-Oliver/go-osc) library under the hood for the osc server.
You can use that library to send messages to VRC.

## Features
-   Auto-discovery for VRC
-   OSC Server
-   OSC Messages
-   Supports the following OSC argument types:
    -   'i' (Int32)
    -   'f' (Float32)
    -   's' (string)
    -   'T' (True)
    -   'F' (False)
    -   'N' (Nil)


## Install

```shell
go get github.com/StarLight-Oliver/vrc_osc_query
```

## Usage

### Server

```go
import "github.com/StarLight-Oliver/vrc_osc_query"

func main() {
	vrcService, err := vrc_osc_query.NewVRCOSCService("Jacket Tracker", 9002)

    if err != nil {
		log.Fatal(err)
	}

	vrcService.AddHandler("/avatar/parameters/jacket", vrc_osc_query.OscTypeInt, "The number for which jacket is shown", func(msg *vrc_osc_query.Message) {
		jacketNum := msg.Arguments[0].(bool)
		fmt.Println("ChatShown", jacketNum)
	})

	err := vrcService.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
```
