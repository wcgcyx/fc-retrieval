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

buildlotus:
	docker build -t wcgcyx/lotusbase -f ./lotus/lotusbase/Dockerfile .
	docker build -t wcgcyx/lotusfull -f ./lotus/lotusfull/Dockerfile .

buildregister:
	docker build -t wcgcyx/fc-retrieval/register -f register/Dockerfile .

builditest:
	docker build -t wcgcyx/fc-retrieval/itest -f itest/Dockerfile .