var=$(go run cmd/helper/main.go)
ret_code=$?
if [ $ret_code == 0 ]; then
    docker run -it --network=shared --env DEVINIT=$var wcgcyx/fc-retrieval/itest
fi
