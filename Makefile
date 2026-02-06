.PHONY: test

test:
	docker compose -p urlshort_test -f docker-compose-test.yml up -d postgres
	docker compose -p urlshort_test -f docker-compose-test.yml up --abort-on-container-exit migrate
	go test ./... -count=1
	docker compose -p urlshort_test -f docker-compose-test.yml down -v