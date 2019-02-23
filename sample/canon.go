package main

import (
    "fmt"
    "../rw"
    "../termutil"
)

func main() {
    fmt.Println("Canonical Mode")
    fmt.Print(">> ")
    input, _ := rw.ReadUntil('q')
    fmt.Println(input)

    fmt.Println("\nUncanonical Mode")
    fmt.Print(">> ")
    var term termutil.Termutil
    term.Init()
    term.SetCanon()

    input, _ = rw.ReadUntil('q')
    fmt.Println()
    fmt.Println(input)

    term.LoadBefore()
}
