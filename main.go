package main

import (
    "encoding/xml"
    "io/ioutil"
    "log"
    "net/http"
    "sync"

    "acronis-server/config"
)

type html struct {
    Body []body `xml:"a"`
}

type body struct {
    Content string `xml:",innerxml"`
}

func startServer() {
    log.Fatal(http.ListenAndServe(":" + config.Port, http.FileServer(http.Dir("./files"))))
}

func getFileNames() ([]string, error) {
    res, err := http.Get(config.URL + ":" + config.Port)
    if err != nil || res.StatusCode != 200 {
        return []string{}, err
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return []string{}, err
    }

    v := html{}
    err = xml.Unmarshal([]byte(body), &v)
    if err != nil {
        return []string{}, err
    }

    result := make([]string, 0, len(v.Body))
    for i := range v.Body {
        result = append(result, v.Body[i].Content)
    }

    return result, nil
}

func main() {
    go startServer()

    files, err := getFileNames()
    if err != nil {
        return
    }

    var wg sync.WaitGroup
    for i := range files {
        wg.Add(1)
        go scanFile(&wg, files[i])
    }

    wg.Wait()

    for i := range fileNames {
        wg.Add(1)
        go downloadFile(&wg, fileNames[i])
    }

    wg.Wait()
}
