package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "io"
    "image"
)

func getFileNames(dirname string) ([]string, error) {
    files, err := ioutil.ReadDir(dirname)
    if err != nil {
        return nil, err
    }
    var filenames []string
    for _, f := range files {
        filenames = append(filenames, f.Name())
    }
    return filenames, nil
}

func decodeImages(files []io.Reader) ([]image.Image, error) {
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

    files, err := getFileNames(dirname)
    if err != nil {
        panic(err)
    }
    fmt.Println(files)


    images, err := decodeImages(files)
    if err != nil {
        panic(err)
    }
    fmt.Println(images)
}
