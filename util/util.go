package util

import (
    "fmt"
    "os"
    "image"
    _"image/png"
    _"image/jpeg"
    "io/ioutil"
    "bytes"

    "github.com/mattn/go-sixel"
    "github.com/nfnt/resize"
)

func GetFileNames(dirname string) ([]string, bool, error) {
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

func DecodeImages(filenames []string) ([]image.Image, [][2]uint , error) {
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

func LoadPages(filenames []string, dirname string, width uint, save bool) ([]*bytes.Buffer, int, error) {
    var pages []*bytes.Buffer
    pagenum := 0
    if save {
        filenames, _, err := GetFileNames(dirname + "/sixel_image")
        if err != nil {
            return nil, 0, err
        }
        pagenum = len(filenames)
        for i := 0; i < pagenum; i++ {
            pages = append(pages, new(bytes.Buffer))
            f, err := os.Open(filenames[i])
            if err != nil {
                return nil, 0, err
            }
            defer f.Close()

            stat, err := f.Stat()
            if err != nil {
                return nil, 0, err
            }
            buf := make([]byte, stat.Size())
            f.Read(buf)
            _, err = pages[i].Write(buf)
            if err != nil {
                return nil, 0, err
            }
        }
    } else {
        images, sizes, err := DecodeImages(filenames)
        if err != nil {
            return nil, 0, err
        }

        pagenum = len(images)
        for i := 0; i < pagenum; i++ {
            pages = append(pages, new(bytes.Buffer))
        }

        path := dirname + "/sixel_image"

        if _, err := os.Stat(path); err != nil {
            if err := os.Mkdir(path, 0777); err != nil {
                return nil, 0, err
            }
        }


        for i, img := range images {
            height := uint((float64(width)/float64(sizes[i][1]))*float64(sizes[i][0]))
            img = resize.Resize(width, height, img, resize.NearestNeighbor)
            sixel.NewEncoder(pages[i]).Encode(img)
            name := fmt.Sprintf("%s/sixel_image/%d.six", dirname, i)
            newfile, err := os.Create(name)
            if err != nil {
                return nil, 0, err
            }
            defer newfile.Close()
            newfile.Write(pages[i].Bytes())

            fmt.Printf("\rSlide loading... %d/%d", i+1, pagenum)
        }
        fmt.Println("")
        fmt.Println("Complete!!")
    }
    return pages, pagenum, nil
}

