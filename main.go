package main

import (
    "encoding/xml"
    "math"
    "net/http"
    "os"
    "sync"

    filesPkg "github.com/Alex-CL/acronis-server/utils"
)

const (
    port = "8080"
)

var (
    mutex sync.Mutex
    index int = math.MaxInt64
    fileNames []string
)

type html struct {
    Body []body `xml:"a"`
}

type body struct {
    Content string `xml:",innerxml"`
}

func startServer() {
    log.Fatal(http.ListenAndServe(":" + port, http.FileServer(http.Dir("./files"))))
}

func getFileNames() ([]string, error) {
    res, err := http.Get("http://localhost:" + port)
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
    if err := filesPkg.CreateTempFiles(); err != nil {
        os.Exit(1)
    }

    defer filesPkg.DeleteTempFiles()
    go startServer()

    files, err := getFileNames()
    if err != nil {
        return
    }

    var wg sync.WaitGroup

    for i := range files {
        wg.Add(1)
        go filesPkg.ScanFile(&wg, files[i])
    }

    wg.Wait()

    for i := range fileNames {
        wg.Add(1)
        go filesPkg.DownloadFile(&wg, fileNames[i])
    }

    wg.Wait()
}
