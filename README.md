Hometube
========

A simple application to download videos using the `youtube-dl` utility.

Baby steps:

- [x] CLI to download videos, sync
- [ ] Separate CLI and server with web API, download videos sync
- [ ] Download async, add web API to query queue status
- [ ] Add simple web interface, without showing downloads in progress
- [ ] Use the original title as filename by default, make filename optional
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

Web API
-------

(draft)

`POST /api/download?url=:url&filename=:filename`

`GET /api/search/downloaded`

`GET /api/search/queue`
