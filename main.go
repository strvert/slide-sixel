package main

import (
    "fmt"
    "flag"
    "io/ioutil"
)

func getFiles(dirname string) []string {
    files, err := ioutil.ReadDir(dirname)
    if err != nil {
        panic(err)
    }

    var filenames []string
    for _, f := range files {
        filenames = append(filenames, f.Name())
    }
    return filenames
}

func main() {
    flag.Parse()
    args := flag.Args()
    dirname := args[0]

    files := getFiles(dirname)

    fmt.Println(files)
}
