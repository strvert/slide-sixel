package main

import (
    "fmt"
    "flag"
    "os"
    "bufio"

    "./termutil"
    "./util"
)

func main() {
    var (
        dirname string
        sixupdate bool
        slidewidth uint
        fullscreen bool
    )
    flag.StringVar(&dirname, "d", "0", "slide directory")
    flag.BoolVar(&sixupdate, "u", false, "update saves for six image")
    flag.UintVar(&slidewidth, "s", 0, "set width for slide width")
    flag.BoolVar(&fullscreen, "f", false, "fullscreen")
    flag.Parse()
    if dirname == "0" {
        fmt.Println(fmt.Errorf("Please set slide directory."))
        os.Exit(1)
    }

    var cs termutil.CtrlSeqs
    width := uint(300)

    if slidewidth != 0 {
        width = uint(slidewidth)
        sixupdate = !true
    } else {
        _, width, _ = cs.GetWindowSize()
    }
    if fullscreen {
        cs.ToggleFullScreen()
        defer cs.ToggleFullScreen()
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
    fmt.Print(string(pages[0].Bytes()))

    var term termutil.Termutil
    term.Init()
    defer term.LoadBefore()
    term.SetCanon()
    term.SetEcho(false)
    reader := bufio.NewReader(os.Stdin)

    currpage := 0
    FOR_LABEL:
    for {
        ch, err := reader.ReadByte()
        if err != nil {
            panic(err)
        }

        switch string(ch) {
            case "q":
                break FOR_LABEL

            case "l":
                if currpage < pagenum-1 {
                    currpage += 1
                }
                fmt.Print("\r")
                fmt.Print(string(pages[currpage].Bytes()))

            case "h":
                if currpage > 0 {
                    currpage -= 1
                }
                fmt.Print("\r")
                fmt.Print(string(pages[currpage].Bytes()))
        }
    }

    cs.ClearScreen()
}
