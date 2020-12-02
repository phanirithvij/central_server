# https://stackoverflow.com/a/64644990/8608146
exe(){
    set -x
    "$@"
    { set +x; } 2>/dev/null
}
exe cd client/web
echo "Building web app"
if [ ! -d node_modules ]; then
    exe npm install
fi
exe npm run build
exe cd ../..
