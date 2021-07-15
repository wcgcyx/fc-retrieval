var=$(go run cmd/helper/main.go)
ret_code=$?
if [ $ret_code == 0 ]; then
    docker run -it --network=shared --env DEV_TEST=$var wcgcyx/fc-retrieval/itest
fi
