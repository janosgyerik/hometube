Hometube
========

A simple application to download videos using the `youtube-dl` utility.

Baby steps:

- [x] CLI to download videos, sync
- [x] Separate CLI and server with web API, download videos sync
- [x] Download async
- [ ] Add simple web interface, without showing downloads in progress
- [ ] Use the original title as filename by default, make filename optional
- [ ] Add web API to query queue status
- [ ] Show download progress (= file size)

Command line interface
----------------------

Examples:

    # short video of 17 seconds
    go run main.go 'https://www.youtube.com/watch?v=H0FcOPb-9rE' blofeld.avi

    # install and run
    go install
    hometube -help
    hometube 'https://www.youtube.com/watch?v=H0FcOPb-9rE' blofeld.avi

    curl -sX POST 'localhost:8080/api/v1/download?url=https://www.youtube.com/watch?v=H0FcOPb-9rE&filename=blofeld.avi' | jq .

Web API
-------

(draft)

`POST /api/v1/download?url=:url&filename=:filename`

`GET /api/v1/search/downloaded`

`GET /api/v1/search/queue`
