# go-dexcomshare

A Golang module that communicates with the Dexcom Share API that was documented in this [Github Gist](https://gist.github.com/StephenBlackWasAlreadyTaken/adb0525344bedade1e25?permalink_comment_id=4001119). 

## Features
- Read most recent estimated glucose value from a Dexcom CGM.
- Retreive the most recent X number of estimated glucose values from a Dexcom CGM.

## Installation

Not installable yet, working on it :)

## Examples
Reading a single estimated glucose value.
```go
package main

import (
    "fmt"
    "log"
    "github.com/cjcocokrisp/go-dexcomshare"
)

func main() {
    dexcom, err := dexcomshare.Login("USERNAME", "PASSWORD")
    if err != nil {
        log.Fatal(err)
    }

    egv, err := dexcom.GetLatestEGV()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(egv.Value)
}
```

Reading multiple glucose values.
```go
package main

import (
    "fmt"
    "log"
    "github.com/cjcocokrisp/go-dexcomshare"
)

func main() {
    dexcom, err := dexcomshare.Login("USERNAME", "PASSWORD")
    if err != nil {
        log.Fatal(err)
    }

    egvs, err := dexcom.GetEGV(5, dexcomshare.DefaultMinutes)
    if err != nil {
        log.Fatal(err)
    }

    // Is retrieved from most recent onwards
    for _, egv := range egvs {
        fmt.Println(egv.Value)
    }
}
```

## Issues & Suggestions

If you have any issues with this module please open an issue on the repo that clearly describes your problem. 

If you have any feature suggestions that you want added also open an issue that describes the feature. If you decide to build it yourself feel free to fork and then open a pull request!