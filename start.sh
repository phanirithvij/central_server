# https://stackoverflow.com/a/64644990/8608146
exe(){
    set -x
    "$@"
    { set +x; } 2>/dev/null
}
exe go generate -x ./...
cd server
exe go build -x
cd ..
exe ./server/server serve $@
