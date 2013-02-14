package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "regexp"
    "strings"
    "strconv"
    "os"
)

const YT_URL = "http://www.youtube.com/watch?v="

// Get Formatlist for given videoId
func getFormatList(videoId string) (map[int]string, error) {
    result := make(map[int]string, 0)
    reg, _ := regexp.Compile("\"url_encoded_fmt_stream_map\": \"([^\"]*)")
    regItag, _ := regexp.Compile("itag=([0-9]+)")
    resp, err := http.Get(fmt.Sprintf("%s%s", YT_URL, videoId))

    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    url_encoded := reg.FindSubmatch(body)

    resultString := string(url_encoded[1])
    resultString = replace(resultString, []string { "%25", "\\u0026", "\\" }, []string { "%", "&", "" })

    // do some replace magic on the url for each video format
    for _, v := range strings.Split(resultString, ",") {

        t := strings.SplitN(v, "url=http", 2)
        v = "url=http" + t[1] + "&" + t[0]
        v = strings.Replace(v, "url=http%3A%2F%2F", "http://", 1)
        v = replace(v, []string { "%3F", "%2F", "%3D", "%26", "%252C", "\\u0026", "sig=" },
                       []string { "?",   "/",   "=",   "&",   "%2C",   "&",       "signature=" })

        itag, _ := strconv.Atoi(regItag.FindStringSubmatch(v)[1])

        if strings.Count(v, "itag=") > 1 {
            v = strings.Replace(v, fmt.Sprintf("&itag=%d", itag), "", 1)
        }

        // Add video url to result
        result[itag] = v
    }

    return result, nil
}

// Download Video from given url
func downloadVideo(url string, filename string) error {
    buffer := make([]byte, 1024)

    fmt.Println("Create file", filename)

    file, err := os.Create(filename)
    if err != nil { return err }
    defer file.Close()

    resp, err := http.Get(url)
    if err != nil { return err }
    defer resp.Body.Close()

    fmt.Printf("Reading %d kB", resp.ContentLength / 1000)

    writeCounter := 0
    for {
        n, err := resp.Body.Read(buffer);

        if n == 0 || err != nil {
            break
        }

        n2, _ := file.Write(buffer[:n])
        writeCounter += n2

        if writeCounter > 1000000 {
            writeCounter = 0
            fmt.Print(".");
        }
    }

    fmt.Println("done")

    return nil
}

// Replace each string in search with
// corresponding string in replace
func replace(value string, search []string, replace []string) string {
    for i := range replace {
        value = strings.Replace(value, search[i], replace[i], -1)
    }

    return value
}

func main() {
    videoId := "8k9XnnqACdc"

    videos, _ := getFormatList(videoId)

    if url, ok := videos[22]; ok {
        downloadVideo(url, videoId + ".mpg")
    }
}
