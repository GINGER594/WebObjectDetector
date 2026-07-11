package main

import (
    "os"
    "fmt"
    "slices"
    "strings"

    argparsing "webobjectdetector/session"
    urlscanning "webobjectdetector/scanning"
)

const helpMessage = "\n\n" +
    "=== Web Object Detector (WOD) ===============================================================================================|\n\n" +
    "basic use:\n" +
    "  ./webobjectdetector <URL> <path-to-words-file>\n\n" +
    "optional flags:\n" +
    "  -ua : path to User-Agents list - one will be selected at random for each GET request\n" +
    "  -p  : pool-size - the maximum number of concurrent GET requests (default 256)\n" +
    "  -t  : timeout threshold in milliseconds (default 30000)\n\n" +
    "output:\n" +
    "  any URLs that return status codes of 400 (generic error, e.g. malformed req or OS error), 404 or 429 are not displayed\n" +
    "  URLs that returned any other status code are shown underneath their status code\n\n" +
    "=============================================================================================================================|\n\n"

func outputScanResults(aggResponses map[int][]string) {
    out := "\n\n===== SCAN COMPLETE =======================================================|\n"
    //dont need to check if map[n] exists becuase the 0-value of a slice is an empty slice (nil):
    out += fmt.Sprintf("[400]: %d\n[404]: %d\n[429]: %d\n\n", len(aggResponses[400]), len(aggResponses[404]), len(aggResponses[429]))
    delete(aggResponses, 400)
    delete(aggResponses, 404)
    delete(aggResponses, 429)
    for statusCode, urls := range aggResponses {
        out += fmt.Sprintf("[%d]: %d\n%s\n\n", statusCode, len(urls), strings.Join(urls, "\n"))
    }
    out += "===========================================================================|\n\n"
    fmt.Println(out)
}

//parses arguments into a settings obj for this session, runs the scan, aggregates & outputs results
func runTool() error {
    if slices.Contains(os.Args, "-help") {
        fmt.Println(helpMessage)
        return nil
    }
    sessionArgs := &argparsing.SessionArgs{}
    err := sessionArgs.ParseArgs(os.Args)
    if err != nil {
        return fmt.Errorf("%s\nrun with the '-help' flag for more info\n", err.Error())
    }

    urlScanner := urlscanning.UrlScanner{}
    urlScanner.ApplySessionArgs(sessionArgs)
    aggResponses := urlScanner.Scan()
    outputScanResults(aggResponses)
    return nil
}

func main() {
    err := runTool()
    if err != nil {
        fmt.Println(err.Error())
    }
}
