up:
	docker compose -f docker/docker-compose.yaml up --build --force-recreate --remove-orphans -d

down:
	docker compose -f docker/docker-compose.yaml down

registry: docker_build docker_save export_images clean 

docker_build:
	docker build -t registry.otlp.example/backend-child ./backend-child
	docker build -t registry.otlp.example/backend-parent ./backend-parent
	docker build -t registry.otlp.example/script-request ./scripts/request

docker_save:
	docker save registry.otlp.example/backend-child > backend-child.tar
	docker save registry.otlp.example/backend-parent > backend-parent.tar
	docker save registry.otlp.example/script-request > script-request.tar

export_images:
	microk8s ctr image import backend-child.tar
	microk8s ctr image import backend-parent.tar
	microk8s ctr image import script-request.tar
	
clean:
	rm backend-child.tar backend-parent.tar script-request.tar
