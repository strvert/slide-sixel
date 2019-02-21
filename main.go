package main

import (
    "fmt"
    "flag"
    "strconv"

    "./termutil"
    "./rw"
    "./util"
)

func main() {
    flag.Parse()
    args := flag.Args()
    dirname := args[0]

    cs := new(termutil.CtrlSeqs)

    _, width, err := cs.GetWindowSize()
    if err != nil {
        panic(err)
    }
    if len(args) >= 2 {
        num, err := strconv.Atoi(args[1])
        if err != nil {
            panic(err)
        }
        width = uint(num)
    }

    filenames, save, err := util.GetFileNames(dirname)
    if err != nil {
        panic(err)
    }

    pages, pagenum, err := util.LoadPages(filenames, dirname, width, save)
    if err != nil {
        panic(err)
    }

    fmt.Println("")
    fmt.Println("Complete!!")

    fmt.Println(string(pages[0].Bytes()))

    currpage := 0
    FOR_LABEL:
    for {
        commands, err := rw.ScanCommand()
        if err != nil {
            panic(err)
        }
        switch commands[0] {
            case "exit", "q":
                break FOR_LABEL

            case "next", "l":
                if currpage < pagenum-1 {
                    currpage += 1
                    fmt.Println(string(pages[currpage].Bytes()))
                }

            case "back", "h":
                if currpage > 0 {
                    currpage -= 1
                    fmt.Println(string(pages[currpage].Bytes()))
                }

            case "jmp", "j":
                page, err := strconv.Atoi(commands[1])
                if err != nil {
                    panic(err)
                }
                page = page-1
                if page < pagenum && page >= 0 {
                    fmt.Println(string(pages[page].Bytes()))
                    currpage = page
                } else {
                    fmt.Println("That page is out range.")
                }
        }
    }
}
