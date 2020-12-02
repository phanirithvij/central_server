# https://stackoverflow.com/a/64644990/8608146
exe(){
    set -x
    "$@"
    { set +x; } 2>/dev/null
}

WEB="false"
PACK="false"

debugInfo () {
  echo "Build web:          $WEB"
  echo "Pack bin assets:    $PACK"
}

buildWeb () {
  exe cd client/web

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
  cd server
  exe go build
  cd ..
}


usage() {
  echo "Usage: $0 [-b web and pack] [-w web only] [-p pack only] [-d debug]" 1>&2;
  exit 1;
}

DEBUG="false"

while getopts "bwp:d" o; do
  case "${o}" in
    b)
      WEB="true"
      PACK="true"
      ;;
    w)
      WEB="true"
      ;;
    p)
      PACK="true"
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
