# Techinical and Implementation Details

> This codebase is the central meta organization managing server for the Corpora collection project by IIIT Hyderabad's LTRC, Speech Processing Lab

We plan on implementing it in golang but golang might be a tough choice as it makes it hard to maintain it. Most people know python but not golang.

- Might use grpc for go/python interop
  - [Example](https://github.com/Jigar3/gRPC_101/blob/master/server/server.go)

Will add in details as I implement the features.

The server uses [cobra](https://github.com/spf13/cobra/) cli tool for command generation

Bash completions also work for linux systems

## TODO
- [ ] Remove pkger and use packr2 because pkger slow

## Setup

- Install go, setup $GOBIN, $GOPATH, \$GOROOT etc..
- `git clone https://github.com/phanirithvij/central_server`
- `bash ./start.sh`
- `bash ./start.sh -d` for debug

Or

```sh
go generate ./...
go build
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
