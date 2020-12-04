# Techinical and Implementation Details

> This codebase is the central meta organization managing server for the Corpora collection project by IIIT Hyderabad's LTRC, Speech Processing Lab

We plan on implementing it in golang but golang might be a tough choice as it makes it hard to maintain it. Most people know python but not golang.

- Might use grpc for go/python interop
  - [Example](https://github.com/Jigar3/gRPC_101/blob/master/server/server.go)

Will add in details as I implement the features.

The server uses [cobra](https://github.com/spf13/cobra/) cli tool for command generation

Bash completions also work for linux systems

## TODO
- [x] Remove pkger and use packr2 because pkger slow
  - [x] Removed packr instead because file modtimes are not supported packr

## Dev notes
- Using client/web a react app might change it to Vue
- React subroute serving [package.json:homepage](https://stackoverflow.com/a/55854101/8608146)
- Use build.sh, start.sh
- `go generate ./server/...` not `go generate ./...` in pwd:central_server
  - Because `central_server` is not package main but `server` is
- `central.go` file should exist if go get needs to work properly
- [filebrowser](https://github.com/phanirithvij/filebrowser) has go rice use along with vue static serving under subroute
  - [fate](https://github.com/phanirithvij/fate) has filebrowser embedded
  - [ ] TODO move it to this
  - [ ] `browser` build tag in `fate` for inbuilt filebrowser support not useful when using s3 or other cloud providers

## Setup

- Install go, setup $GOBIN, $GOPATH, \$GOROOT etc..
- `git clone https://github.com/phanirithvij/central_server`
- `bash ./start.sh`
- `bash ./start.sh -d` for debug

Or

```sh
# client
cd client/react
npm i
npm run build
cd ../..
# server
cd server
go generate ./...
go build
cd ..
./server/server serve # -d
```

## Routes

```
/home
/register
/api
/api/v1
/api/v2
```

---

```
/dashboard
/privacy
/terms
/docs
/docs/api
/docs/orgs
```

## Server plan

- Home page which has some info about what it is
- Privacy Policy page
- Terms of service page
- Register page
  - Register for organizations
  - Then they get redirected to the docs which are filled in with their fields and setup etc.
- Central Dashboard
  - Contains a settings page
  - A server monitoring page
  - Manage Organizations
    - Manual add new (from requests)
    - Delete org (from reports)
    - Disable org (temp delete)
    - Details page for each org
      - Org analytics
      - User stats (forwarded if public)
      - User reports
      - Support requests
  - Global Analytics
- Organization console
  - Delete org (self request permanet)
  - Unlink org (self request temp delete)
  - Privacy settings
  - Corpora data License
- Documentation
  - Swagger (?)
    - API for clients
    - Organization setup docs
  - Readthedocs for general docs (?)

## Client Pages

- Settings
  - User Privacy
  - Linked Organizations
  - Usage mode
    - Guest
    - Login via different providers (default)
- File upload with https://github.com/jjmutumi/tus_client.git
