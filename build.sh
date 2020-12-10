# https://stackoverflow.com/a/64644990/8608146
exe(){
    set -x
    "$@"
    { set +x; } 2>/dev/null
}

WEB="false"
PACK="false"
BIN="false"

EXT=""

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN|MINGW|MSYS*)
        machine=Windows
        EXT=".exe"
    ;;
    *)          machine="UNKNOWN:${unameOut}"
esac

debugInfo () {
    echo "Build web:          $WEB"
    echo "Pack bin assets:    $PACK"
    echo "Build binary:       $BIN"
}

buildWebAdmin () {
    exe cd client/admin
    
    if [ ! -d "node_modules" ]; then
        if [ "$CI" = "true" ]; then
            exe npm ci
        else
            exe npm install
        fi
    fi
    
    exe npm run build
    exe cd ../..
}


buildWebOrg () {
    exe cd client/org
    
    if [ ! -d "node_modules" ]; then
        if [ "$CI" = "true" ]; then
            exe npm ci
        else
            exe npm install
        fi
    fi
    
    exe npm run build
    exe cd ../..
}

cleanPacked() {
    echo "Removing generated go bin files..."
    # need to install rimraf for cross platform rm -rf glob
    if ! [ -x "$(command -v rimraf)" ]
    then
        npm i -g rimraf
    fi
    
    cd server
    exe rimraf -g "**/*_g.go"
    cd ..
}

packAssets () {
    echo "Packing assets..."
    if ! [ -x "$(command -v pkger)" ]
    then
        exe go get -u -v github.com/markbates/pkger/cmd/pkger
    fi
    exe go generate -x ./server/...
}

buildBin () {
    echo "Building binary..."
    exe cd server
    exe go build -o server-prod$EXT &
    exe go build -o wsserver-prod$EXT -ldflags "-w -s"
    wait
    exe cd ..
}

buildBinDev () {
    echo "Building binary for development..."
    exe cd server
    exe go build -tags skippkger &
    # exe go build -tags skippkger -o ws-server$EXT -ldflags "-w -s"
    wait
    exe cd ..
}

usage() {
    echo "Usage: $0 [-a web,pack,build] [-w web only] [-p pack only] [-b build only] [-d dev build] [-c clean _g.go] [-v verbose]" 1>&2;
    exit 1;
}

DEBUG="false"

while getopts "acwpbd:v" o; do
    case "${o}" in
        a)
            WEB="true"
            PACK="true"
            BIN="true"
        ;;
        c)
            CLEAN="true"
        ;;
        w)
            WEB="true"
        ;;
        p)
            PACK="true"
        ;;
        b)
            BIN="true"
        ;;
        d)
            DEV="true"
            BIN="false"
        ;;
        v)
            DEBUG="true"
        ;;
        *)
            usage
        ;;
    esac
done
shift $((OPTIND-1))

if [ "$DEBUG" = "true" ]; then
    debugInfo
fi

if [ "$CLEAN" = "true" ]; then
    cleanPacked
fi

if [ "$WEB" = "true" ]; then
    buildWebAdmin &
    buildWebOrg
fi

if [ "$PACK" = "true" ]; then
    wait
    packAssets
fi

if [ "$DEV" = "true" ]; then
    buildBinDev
fi

if [ "$BIN" = "true" ]; then
    buildBin
fi
