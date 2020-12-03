# https://stackoverflow.com/a/64644990/8608146
exe(){
    set -x
    "$@"
    { set +x; } 2>/dev/null
}
if [ ! -f "server/server" ] && [ ! -f "server/server.exe" ]; then
    echo "Make sure to run"
    echo -e "\tsh build.sh -b -d"
    echo "before running start.sh"
    exit 1
fi
exe ./server/server serve $@
