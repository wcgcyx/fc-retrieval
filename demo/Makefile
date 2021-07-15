start:
	docker network create shared || true
	docker stop $(shell docker ps -q) || true
	docker rm $(shell docker ps -q -a) || true
	docker compose up

clean:
	docker stop $(shell docker ps -q) || true
	docker rm $(shell docker ps -q -a)