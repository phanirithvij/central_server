# https://stackoverflow.com/a/64644990/8608146
exe(){
    set -x
    "$@"
    { set +x; } 2>/dev/null
}
if ! command -v packr2
then
    exe go get -u -v github.com/gobuffalo/packr/v2/packr2
fi
exe go generate ./...
cd server
# exe go build
exe go build -x
cd ..
exe ./server/server serve $@
