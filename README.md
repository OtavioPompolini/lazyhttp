# Project Postman

## Commands
    n -> create new request
    enter -> select request to edit
    P -> perform the request
    D -> delete selected request
    j -> navigate down through requests
    k -> navigate up through requests

    Response window:
        V -> Visual line mode
        y -> Copy content selected in visual mode

    *Capital letters = Shift + <key>

#### REQUEST SINTAX
    <METHOD> <URL>
    <HeaderKey>=<headerValue>
    <HeaderKey>=<headerValue>
    <HeaderKey>=<headerValue>

    <JsonBody>

*It needs an empty line between headers and body

## OBS. v0.2
Migrating from v0.1 will not carry database through v0.2

## Known issues
    - Not copying to clipboard on linux

## Roadmap some day
    - Response headers not always display in same order
    - UI improvements based on OS window size (Windows rescale)
    - DB Migrations
    - FIX THE GOD DAMN DLL (I was too dumb to do it correctly)
    - Able to write path param in different lines
    - find/filter request
    - Request navigation Ctrl+u, Crtl+d -> jump multiple
    - Windows navigation through hjkl
    - Add support to collections
    - Add support to environments and variables
    - Add VIMOTION when editing request (THIS IS A MUST)
    - ? keybind to open a help menu
    - Support other formats url-encoded
    - Delete request confirmation
    - Save all on close
