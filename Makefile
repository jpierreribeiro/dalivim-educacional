build:
	docker compose up --build -d
buildapi:
	 docker compose up --build api
db:
	docker exec -it dalivim-educacional psql -U p -d dalivim
docker:
	docker start dalivim-educacional
posting:
	posting --collection posting

.PHONY: build, buildapi, db, posting, docker