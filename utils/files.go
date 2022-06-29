package utils

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
)

type fileContent struct {
    file string
    text string
}

var contents = []fileContent{
    fileContent{ text: "---A---" },
    fileContent{ text: "--A------" },
    fileContent{ text: "------------" },
    fileContent{ text: "==A==========" },
}

func scanFile(wg *sync.WaitGroup, file string) {
    defer wg.Done()

    res, err := http.Get(url + "/" + file)
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

    resp, err := http.Get(url + "/" + file)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    out, err := os.Create("./" + file)
    if err != nil {
        return
    }
    defer out.Close()

    if _, err = io.Copy(out, resp.Body); err != nil {
        return
    }
}

func createTempFiles() error {
    for i := range contents {
    	f, err := os.CreateTemp("", fmt.Sprintf("file%d", i + 1))
        if err != nil {
    		return err
    	}
        contents[i].file = f.Name()

        if _, err := f.Write([]byte(contents[i].text)); err != nil {
    		f.Close()
    		return err
    	}
    	if err := f.Close(); err != nil {
    		return err
    	}
    }

    return nil
}

func deleteTempFiles() error {
    for i := range contents {
    	if err := os.Remove(contents[i].file); err != nil {
    		return err
    	}
    }

    return nil
}
