.PHONY: githooks
githooks:
	git config --local core.hooksPath .githooks/

.PHONY: loc
loc:
	cloc --3 --exclude-list-file=.clocignore .

.PHONY: crq
crq:
	docker exec -it talvi-cockroach-1 cockroach sql --insecure