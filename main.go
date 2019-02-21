package main

import (
    "fmt"
    "flag"
    "strconv"
    "os"

    "./termutil"
    "./rw"
    "./util"
)

func main() {
    var (
        dirname string
        sixupdate bool
        slidewidth uint
    )
    flag.StringVar(&dirname, "d", "0", "slide directory")
    flag.BoolVar(&sixupdate, "u", false, "update saves for six image")
    flag.UintVar(&slidewidth, "s", 0, "set width for slide width")
    flag.Parse()
    if dirname == "0" {
        fmt.Println(fmt.Errorf("Please set slide directory."))
        os.Exit(1)
    }

    var cs termutil.CtrlSeqs

    _, width, err := cs.GetWindowSize()
    if err != nil {
        panic(err)
    }
    if slidewidth != 0 {
        width = uint(slidewidth)
        sixupdate := !true
    }

    filenames, save, err := util.GetFileNames(dirname)
    if err != nil {
        panic(err)
    }

    save = save && !sixupdate
    pages, pagenum, err := util.LoadPages(filenames, dirname, width, save)
    if err != nil {
        panic(err)
    }

    cs.ClearScreen()
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

    cs.ClearScreen()
}
