package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "os"
    "bytes"
    "math"
    "strconv"
    "image"
    _"image/png"

    "github.com/mattn/go-sixel"
    "github.com/nfnt/resize"

    "./termutil"
    "./rw"
)

func getFileNames(dirname string) ([]string, error) {
    files, err := ioutil.ReadDir(dirname)
    if err != nil {
        return nil, err
    }
    var filenames []string
    for _, f := range files {
        path := fmt.Sprintf("%s/%s", dirname, f.Name())
        filenames = append(filenames, path)
    }
    return filenames, nil
}

func decodeImages(filenames []string) ([]image.Image, error) {
    var files []*os.File
    for _, f := range filenames {
        file, err := os.Open(f)
        if err != nil {
            return nil, err
        }
        files = append(files, file)
    }
    var images []image.Image
    for _, f := range files {
        img, _, err := image.Decode(f)
        if err != nil {
            return nil, err
        }
        images = append(images, img)
    }
    return images, nil
}

func main() {
    flag.Parse()
    args := flag.Args()
    dirname := args[0]

    maxwidth := uint(400)
    if len(args) >= 2 {
        num, err := strconv.Atoi(args[1])
        if err != nil {
            panic(err)
        }
        maxwidth = uint(num)
    }

    files, err := getFileNames(dirname)
    if err != nil {
        panic(err)
    }

    images, err := decodeImages(files)
    if err != nil {
        panic(err)
    }
    pagenum := len(images)

    var writer []*bytes.Buffer
    for i := 0; i < pagenum; i++ {
        writer = append(writer, new(bytes.Buffer))
    }

    for i, img := range images {
        img = resize.Thumbnail(maxwidth, math.MaxUint32, img, resize.NearestNeighbor)
        sixel.NewEncoder(writer[i]).Encode(img)
        fmt.Printf("\rSlide loading... %d/%d", i+1, pagenum)
    }
    fmt.Println("")
    fmt.Println("Complete!!")

    fmt.Println(string(writer[0].Bytes()))

    var term termutil.Termutil
    term.Init()
    term.SetCanon()
    defer term.SetUncanon()

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
                    fmt.Println(string(writer[currpage].Bytes()))
                }

            case "back", "h":
                if currpage > 0 {
                    currpage -= 1
                    fmt.Println(string(writer[currpage].Bytes()))
                }

            case "jmp", "j":
                page, err := strconv.Atoi(commands[1])
                if err != nil {
                    panic(err)
                }
                page = page-1
                if page < pagenum && page >= 0 {
                    fmt.Println(string(writer[page].Bytes()))
                    currpage = page
                } else {
                    fmt.Println("That page is out range.")
                }
        }
    }
}
