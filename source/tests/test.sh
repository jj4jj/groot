if [ $# -eq 1 ];then
    set -x
    go test -v -run Test $1_test.go
    set +x
elif [ $# -gt 1 ];then
    set -x
    go test -v -run $2 $1_test.go
    set +x
fi

