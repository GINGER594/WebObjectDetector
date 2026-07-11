package urlscanning

import (
    "fmt"
    "net/http"
    "time"
    "math/rand"
    "strings"

    argparsing "webobjectdetector/session"
)

const progressBarLen int = 32

//hashmap of reserved IPv4 addresses in range [0:224] (beyond 224 is generally reserved)
var reservedIPv4Adrs map[int]bool = map[int]bool{
    0:true,
    10:true,
    100:true,
    127:true,
    169:true,
    172:true,
    192:true,
    198:true,
    203:true,
}

//struct for storing the statusCode for an httpResponse with the url to id it
type httpResponse struct {
    url string
    statusCode int
}

type UrlScanner struct {
    settings *argparsing.SessionArgs
}

func aggregateResponses(responses []httpResponse) map[int][]string {
    aggResponses := map[int][]string{}
    for _, resp := range responses {
        if storedResponses, ok := aggResponses[resp.statusCode]; !ok {
            aggResponses[resp.statusCode] = []string{resp.url}
        } else {
            aggResponses[resp.statusCode] = append(storedResponses, resp.url)
        }
    }
    return aggResponses
}

func showProgressBar(got, expected int) {
    progress := int(float64(progressBarLen) * (float64(got) / float64(expected)))
    progressBar := strings.Repeat("#", progress) + strings.Repeat(" ", progressBarLen-progress)
    fmt.Printf("\x1b[0G\x1b[0Kscanning: [%s] [%d/%d] \x1b[0K", progressBar, got, expected)
}

func getRandUserAgent(userAgents []string) string {
    length := len(userAgents)
    randIndex := rand.Intn(length)
    return userAgents[randIndex]
}

//generates a random valid IPv4 add
func genRandIPv4() string {
    var first int
    for {
        first = rand.Intn(224)
        if _, ok := reservedIPv4Adrs[first]; !ok {
            break
        }
    }
    return fmt.Sprintf("%d.%d.%d.%d", first, rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

//sends a GET req to a given url and returns the status code in an httpResponse obj
func (us *UrlScanner) scanUrl(client *http.Client, url string) httpResponse {
    //req set-up
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return httpResponse{url: url, statusCode: 400}
    }
    req.Header.Add("X-Forwarded-For", genRandIPv4())
    if us.settings.UserAgents != nil {
        req.Header.Add("User-Agent", getRandUserAgent(us.settings.UserAgents))
    }

    //sending request
    resp, err := client.Do(req)
    if err != nil {
        return httpResponse{url: url, statusCode: 400}
    }
    defer resp.Body.Close()
    return httpResponse{url: url, statusCode: resp.StatusCode}
}

//wrapper func designed to be run asynchronously, blocks if grPool is full
func (us *UrlScanner) asyncScanUrl(grPool chan bool, responseCh chan httpResponse, client *http.Client, url string) {
    grPool <- true
    responseCh <- us.scanUrl(client, url)
    <-grPool
}

//performs a full scan according to the session args and returns the results
func (us *UrlScanner) Scan() map[int][]string {
    //client set-up
    client := &http.Client{
        Timeout: time.Duration(us.settings.Timeout) * time.Millisecond,
    }
    defer client.CloseIdleConnections()

    var responses []httpResponse
    grPool := make(chan bool, us.settings.PoolSize) //buffered channel goroutine-pool, sending/receiving blocks when full
    responseCh := make(chan httpResponse)
    for _, word := range us.settings.Words {
        url := us.settings.BaseUrl + "/" + word
        go us.asyncScanUrl(grPool, responseCh, client, url)        
    }
    for ; len(responses) < len(us.settings.Words); {
        responses = append(responses, <-responseCh)
        showProgressBar(len(responses), len(us.settings.Words))
    }

    return aggregateResponses(responses)
}

//takes in sessionArgs, applies to settings, constructs http client
func (us *UrlScanner) ApplySessionArgs(sessionArgs *argparsing.SessionArgs) {
    //applying default values if they have not been provided by the user
    if sessionArgs.PoolSize <= 0 {
        sessionArgs.PoolSize = argparsing.DefaultPoolSize
    }
    if sessionArgs.Timeout <= 0 {
        sessionArgs.Timeout = argparsing.DefaultTimeout
    }

    us.settings = sessionArgs
}
