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
    _"image/jpeg"

    "github.com/mattn/go-sixel"
    "github.com/nfnt/resize"

    "./termutil"
    "./rw"
)

func getFileNames(dirname string) ([]string, bool, error) {
    files, err := ioutil.ReadDir(dirname)
    save := false
    if err != nil {
        return nil, false, err
    }
    var filenames []string
    for _, f := range files {
        name := f.Name()
        if f.IsDir() {
            if name == "sixel_image" {
                save = true
            }
        } else {
            path := fmt.Sprintf("%s/%s", dirname, name)
            filenames = append(filenames, path)
        }
    }
    return filenames, save, nil
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

    filenames, save, err := getFileNames(dirname)
    if err != nil {
        panic(err)
    }

    var pages []*bytes.Buffer
    pagenum := 0
    if save {
        filenames, save, err = getFileNames(dirname + "/sixel_image")
        if err != nil {
            panic(err)
        }
        pagenum = len(filenames)
        for i := 0; i < pagenum; i++ {
            pages = append(pages, new(bytes.Buffer))
            f, err := os.Open(filenames[i])
            if err != nil {
                panic(err)
            }
            defer f.Close()

            stat, err := f.Stat()
            if err != nil {
                panic(err)
            }
            buf := make([]byte, stat.Size())
            f.Read(buf)
            _, err = pages[i].Write(buf)
            if err != nil {
                panic(err)
            }
        }
    } else {
        images, sizes, err := decodeImages(filenames)
        if err != nil {
            panic(err)
        }

        pagenum = len(images)
        for i := 0; i < pagenum; i++ {
            pages = append(pages, new(bytes.Buffer))
        }

        if err := os.Mkdir(dirname + "/sixel_image", 0777); err != nil {
            panic(err)
        }

        for i, img := range images {
            height := uint((float64(width)/float64(sizes[i][1]))*float64(sizes[i][0]))
            img = resize.Resize(width, height, img, resize.NearestNeighbor)
            sixel.NewEncoder(pages[i]).Encode(img)
            name := fmt.Sprintf("%s/sixel_image/%d.six", dirname, i)
            newfile, err := os.Create(name)
            if err != nil {
                panic(err)
            }
            defer newfile.Close()
            newfile.Write(pages[i].Bytes())

            fmt.Printf("\rSlide loading... %d/%d", i+1, pagenum)
        }
    }
    fmt.Println(filenames)
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
