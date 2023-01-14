up:
	docker compose -f docker/docker-compose.yaml up --build --force-recreate --remove-orphans -d

down:
	docker compose -f docker/docker-compose.yaml down 