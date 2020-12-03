# https://stackoverflow.com/a/64644990/8608146
exe(){
    set -x
    "$@"
    { set +x; } 2>/dev/null
}

WEB="false"
PACK="false"
BIN="false"

debugInfo () {
  echo "Build web:          $WEB"
  echo "Pack bin assets:    $PACK"
  echo "Build web:          $BIN"
}

buildWeb () {
  exe cd client/vue
  # exe cd client/web

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

packAssets () {
  if ! [ -x "$(command -v packr2)" ]
  then
    exe go get -u -v github.com/gobuffalo/packr/v2/packr2
  fi
  exe go generate -x ./server/...
}

buildBin () {
  exe cd server
  exe go build
  exe cd ..
}

usage() {
  echo "Usage: $0 [-a web,pack,build] [-w web only] [-p pack only] [-b build only] [-d debug]" 1>&2;
  exit 1;
}

DEBUG="false"

while getopts "awpb:d" o; do
  case "${o}" in
    a)
      WEB="true"
      PACK="true"
      BIN="true"
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

if [ "$WEB" = "true" ]; then
  buildWeb
fi

if [ "$PACK" = "true" ]; then
  packAssets
fi

if [ "$BIN" = "true" ]; then
  buildBin
fi
