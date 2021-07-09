build:
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