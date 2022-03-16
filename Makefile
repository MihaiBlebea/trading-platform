docker-build:
	docker build -t serbanblebea/trading-platform:v1.0 .

docker-run:
	docker run --env-file ./.env -d -p 8080:8080 --name trading_platform serbanblebea/trading-platform:v1.0

docker: docker-build docker-run

docker-stop:
	docker stop trading_platform && docker rm trading_platform