package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "os"
    "bytes"
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

func decodeImages(filenames []string) ([]image.Image, [][2]uint , error) {
    var files []*os.File
    var images []image.Image
    var sizes [][2]uint
    for _, f := range filenames {
        file, err := os.Open(f)
        if err != nil {
            return nil, [][2]uint{{0, 0}}, err
        }
        defer file.Close()
        files = append(files, file)

        img, _, err := image.Decode(file)
        if err != nil {
            return nil, [][2]uint{{0, 0}}, err
        }
        images = append(images, img)

        file, err = os.Open(f)
        if err != nil {
            return nil, [][2]uint{{0, 0}}, err
        }
        defer file.Close()

        imgconf, _, err := image.DecodeConfig(file)
        if err != nil {
            return nil, [][2]uint{{0, 0}}, err
        }
        sizes = append(sizes, [2]uint{uint(imgconf.Height), uint(imgconf.Width)})
    }
    return images, sizes, nil
}

func main() {
    var term termutil.Termutil

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

    filenames, err := getFileNames(dirname)
    if err != nil {
        panic(err)
    }

    images, sizes, err := decodeImages(filenames)
    if err != nil {
        panic(err)
    }
    pagenum := len(images)

    var writer []*bytes.Buffer
    for i := 0; i < pagenum; i++ {
        writer = append(writer, new(bytes.Buffer))
    }

    for i, img := range images {
        height := uint((float64(width)/float64(sizes[i][1]))*float64(sizes[i][0]))
        img = resize.Resize(width, height, img, resize.NearestNeighbor)
        sixel.NewEncoder(writer[i]).Encode(img)
        fmt.Printf("\rSlide loading... %d/%d", i+1, pagenum)
    }
    fmt.Println("")
    fmt.Println("Complete!!")

    fmt.Println(string(writer[0].Bytes()))

    term.Init()
    if err := term.SetCanon(); err != nil {
        panic(err)
    }

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

    if err := term.SetUncanon(); err != nil {
        panic(err)
    }
}
