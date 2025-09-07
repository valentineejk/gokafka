start:
	docker compose up -d

down:
	docker compose down

run-producer:
	cd ./producer && go run .

run-consumer:
	cd ./consumer && go run .

