buildgw:
	docker build -t wcgcyx/fc-retrieval/gateway -f gateway/Dockerfile .

buildpvd:
	docker build -t wcgcyx/fc-retrieval/provider -f provider/Dockerfile .

buildclient:
	docker build -t wcgcyx/fc-retrieval/client -f client/Dockerfile .

buildgwadmin:
	docker build -t wcgcyx/fc-retrieval/gateway-admin -f gateway-admin/Dockerfile .

buildpvdadmin:
	docker build -t wcgcyx/fc-retrieval/provider-admin -f provider-admin/Dockerfile .

buildall:
	docker build -t wcgcyx/fc-retrieval/gateway -f gateway/Dockerfile .
	docker build -t wcgcyx/fc-retrieval/provider -f provider/Dockerfile .
	docker build -t wcgcyx/fc-retrieval/client -f client/Dockerfile .
	docker build -t wcgcyx/fc-retrieval/gateway-admin -f gateway-admin/Dockerfile .
	docker build -t wcgcyx/fc-retrieval/provider-admin -f provider-admin/Dockerfile .
	docker build -t wcgcyx/fc-retrieval/register -f register/Dockerfile .
	docker build -t wcgcyx/lotusbase -f lotusbase/Dockerfile .
	docker build -t wcgcyx/lotusfull -f lotusfull/Dockerfile .

start:
	docker compose up

clean:
	docker stop $(shell docker ps -q) || true
	docker rm $(shell docker ps -q -a)

gen:
	cd util; go run ./main.go; cd ..

runclient:
	docker run -it --network=shared wcgcyx/fc-retrieval/client ./main

rungwadmin:
	docker run -it --network=shared wcgcyx/fc-retrieval/gateway-admin ./main

runpvdadmin:
	docker run -it --network=shared wcgcyx/fc-retrieval/provider-admin ./main