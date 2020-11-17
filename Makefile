run-services:
	docker stack deploy --compose-file deployments/services.yml gen
rm:
	docker service rm $(shell docker service ls -q) || true
