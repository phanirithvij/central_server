# Techinical and Implementation Details

> This codebase is the central meta organization managing server for the Corpora collection project by IIIT Hyderabad's LTRC, Speech Processing Lab

We plan on implementing it in golang but golang might be a tough choice as it makes it hard to maintain it. Most people know python but not golang.

- Might use grpc for go/python interop
  - [Example](https://github.com/Jigar3/gRPC_101/blob/master/server/server.go)

Will add in details as I implement the features.

The server uses [cobra](https://github.com/spf13/cobra/) cli tool for command generation

Bash completions also work for linux systems

**No typescript** -> slightly more errors but huge time waste handling typescript issues

## TODO

- [x] Remove pkger and use packr2 because pkger slow
  - [x] Removed packr instead because file modtimes are not supported packr
- [ ] Use go rice instead to make building faster?
- [ ] TIMEPASS: Add source zip to binary and allow download via endpoint
- [x] [Leaflet-geosearch](https://github.com/smeijer/leaflet-geosearch) for getting location then drag [marker](https://stackoverflow.com/questions/27271994/leaflet-draggable-marker-and-coordinates-display-in-a-field-form) to get latlong
- [ ] Might be useful [timezone-builder](https://github.com/evansiroky/timezone-boundary-builder)
- [ ] Attribute https://icons8.com/icons somewhere
  - Using ok.svg, no.svg
- [ ] Move register endpoints to /api/v1/register
- [ ] Aliases cannot be valid emails write that constraint when validating
  - Because if so then the login method will fail if alias is detected as an email
  - Someone can accidentally create an alias which is someone elses' email and we will cry that day (so remove this bad design)

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
- create-react-app [issue](https://github.com/facebook/create-react-app/issues/1070) with react-scripts start and [workaround](https://github.com/facebook/create-react-app/issues/1070#issuecomment-721477819) is [cra-build-watch](https://github.com/Nargonath/cra-build-watch)
  - This has errors in console from `webpackdevserver.js`

## Setup

- Install go, setup $GOBIN, $GOPATH, \$GOROOT etc..
- `git clone https://github.com/phanirithvij/central_server`
- `bash ./start.sh` or `sh ./start.sh`
- `bash ./start.sh -d` for debug

Or

```sh
# clients
cd client/react
npm i
npm run build
cd ../..
cd client/admin
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

## Development

Client

```sh
cd client/org
# npm install
# for developing along with the server
npm start
# OR to work on just the react app individually
npm start:dev
# go to http://localhost:3000/org
# code
```

```sh
cd client/admin
# npm install
# for developing along with the server
npm start
# OR to work on just the react app individually
npm start:dev
# go to http://localhost:3000/admin
# code
```

Server

```sh
sh build.sh -b -d
sh start.sh -d
```

Go to [http://localhost:9090/admin](http://localhost:9090/admin)

## Deploy

- [ ] TODO use go releaser (??)

```sh
sh build.sh -a -d
# -a will build client/admin and client/org
# then pack assets
# then go build
```

## Routes

Public views

```
/privacy
/terms
/org                    - Default homepage for public (REACT)
  /register               - Organization register & login
  /dashboard          - Organization Dashboard
    /activity             - Organization activity
/docs                   - Documentation
  /docs/api             - Public API documentation
  /docs/orgs            - Public organization documentation
/api                    - Public API lists
  /api/v1               - Public API v1
  /api/v2               - ...
  ...
```

---

Admin views

```
/admin            Admin dashboard (REACT app)
  /activity       Admin activity & logs
  /activity/logs  View & export recent logs
  /analytics      Central analytics
  /list           List of all admins
  /register       Admin registration & login
  /status         Metrics and stuff
  /new            Add a new admin
  /added          Admins added by this admin
  /profile        Admin's profile
  /logout         Admin logout
  /files          View organizations as a filesystem
```

Status.hub.org

```
/status           Status of the server
/incident/id      Incident details
```

- [ssm](https://github.com/ssimunic/gossm) server status notifications example
  - Email, Slack etc. supported

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
  - Swagger (??)
    - API for clients
    - Organization setup docs
  - Readthedocs for general docs (??)
  - Or docz [this](https://github.com/doczjs/docz/) ??

## Client Pages

- Settings
  - User Privacy
  - Linked Organizations
  - Usage mode
    - Guest
    - Login via different providers (default)
- File upload with https://github.com/jjmutumi/tus_client.git
