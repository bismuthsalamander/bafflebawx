package main

import (
    "fmt"
    "github.com/bismuthsalamander/bafflebawx/inceptor"
)
func main() {
    n, err := inceptor.Uint64()
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
    } else {
        fmt.Printf("%v\n", n)
    }
    inceptor.Server().Run()
}
