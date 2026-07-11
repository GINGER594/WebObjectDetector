# WebObjectDetector
An asynchronous terminal tool for finding hidden web-objects

Features:
- asynchronous requesting
- randomised X-Forwarded-For & User-Agent headers
- configurable request-pool size
- configurable timeout threshold

--- How to use ---
Basic use:
- ./wod_YOUR_OS <URL> <path-to-words-file>

Optional flags:
- -ua : path to User-Agents list - one will be selected at random for each GET request
- -p  : pool-size - the maximum number of concurrent GET requests (default 256)
- -t  : timeout threshold in milliseconds (default 30000)

Output:
- any URLs that return status codes of 400 (generic error, e.g. malformed req or OS error), 404 or 429 are not displayed
- URLs that returned any other status code are shown underneath their status code
