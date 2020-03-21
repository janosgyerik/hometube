Hometube
========

A simple application to download videos using the `youtube-dl` utility.

Baby steps:

- [x] CLI to download videos, sync
- [x] Separate CLI and server with web API, download videos sync
- [x] Download async
- [x] Add simple web interface, without showing downloads in progress
- [ ] Clean up and polish
    - refuse to download if filename exists (beware of .mp4 added by youtube-dl!)
    - submit form on Enter in any field
    - add link to open file, or open basedir in file explorer
    - filename should not contain extension
    - fix all TODOs
- [ ] Use the original title as filename by default, make filename optional
- [ ] Add queue content to UI, and periodically query it and update view
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

User interface: visit `/home` in a browser.

Request to download a URL and save as a filename asynchronously:

    POST /api/v1/download?url=:url&filename=:filename

Get the list of downloaded files:

    GET /api/v1/list/downloaded

Get the list of queued downloads (not implemented yet):

    GET /api/v1/list/queue
