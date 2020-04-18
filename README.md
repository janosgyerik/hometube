Hometube
========

A simple application to download videos using the `youtube-dl` utility.

Baby steps:

- [x] CLI to download videos, sync
- [x] Separate CLI and server with web API, download videos sync
- [x] Download async
- [x] Add simple web interface, without showing downloads in progress
- [ ] Make API url configurable for app package
- [ ] Clean up and polish
    - add link to open file, or open basedir in file explorer
    - fix all TODOs
- [x] Use the original title as filename by default, make filename optional
- [x] Support downloading playlists
- [ ] Add queue content to UI, and periodically query it and update view
- [ ] Deploy in the cloud

Command line interface
----------------------

Examples:

    # short video of 17 seconds
    go run main.go 'https://www.youtube.com/watch?v=H0FcOPb-9rE'

    # install and run
    go install
    hometube -help
    hometube 'https://www.youtube.com/watch?v=H0FcOPb-9rE'

    curl -sX POST 'localhost:8080/api/v1/download?url=https://www.youtube.com/watch?v=H0FcOPb-9rE' | jq .

    # download playlist
    curl -sX POST 'localhost:8080/api/v1/download?url=https://www.youtube.com/watch?list=PLvvIqOmW1pWwxMz1cfdJOOC-sMYqteq0a' | jq .

Web API
-------

User interface: visit `/home` in a browser.

Request to download a URL and save as a filename asynchronously:

    POST /api/v1/download?url=:url&filename=:filename

Get the list of downloaded files:

    GET /api/v1/list/downloaded

Get the list of queued downloads (not implemented yet):

    GET /api/v1/list/queue
