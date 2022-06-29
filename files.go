package main

import (
    "io"
    "io/ioutil"
    "math"
    "net/http"
    "os"
    "sync"

    "acronis-server/config"
)

var (
    mutex sync.Mutex
    index int = math.MaxInt64
    fileNames []string
)

func scanFile(wg *sync.WaitGroup, file string) {
    defer wg.Done()

    res, err := http.Get(config.URL + ":" + config.Port + "/" + file)
    if err != nil || res.StatusCode != 200 {
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return
    }

    exit := false
    for i := 0; i < len(body) && !exit; i++ {
        if (body[i] != 'A') {
            continue;
        }

        mutex.Lock()
        if (i < index) {
            index = i
            fileNames = []string{file}
        } else if (i == index) {
            fileNames = append(fileNames, file)
        } else {
            exit = true
        }
        mutex.Unlock()
    }
}

func downloadFile(wg *sync.WaitGroup, file string) {
    defer wg.Done()

    resp, err := http.Get(config.URL + ":" + config.Port + "/" + file)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    out, err := os.CreateTemp("", file)
    if err != nil {
        return
    }
    defer out.Close()

    if _, err = io.Copy(out, resp.Body); err != nil {
        return
    }
}
