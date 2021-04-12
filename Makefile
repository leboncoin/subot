# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# DOCKER TASKS

# Build and run the container
start: ## Spin up the project
	docker-compose up -d

start-build: ## Spin up the project after rebuild
	docker-compose up -d --build

stop: ## Stop running containers
	docker-compose down

logs: ## Show all logs
	docker-compose logs -f

log-analytics: ## Show analytics service logs
	docker-compose logs -f analytics

log-replier: ## Show replier service logs
	docker-compose logs -f replier

restart-replier: ## Restart replier service
	docker-compose restart replier

restart-analytics: ## Restart analytics service
	docker-compose restart analytics

rebuild-replier: ## Rebuild replier service
	docker-compose up -d --build replier

rebuild-analytics: ## Rebuild analytics service
	docker-compose up -d --build analytics

# Variable for filename for store running procees id
PID_FILE = /tmp/my-app.pid
# We can use such syntax to get main.go and other root Go files.
GO_FILES = $(wildcard *.go)

# Start task performs "go run main.go" command and writes it's process id to PID_FILE.
startapp:
	go run /go/src/github.com/leboncoin/subot/services/$(APP)/cmd & echo $$! > $(PID_FILE)
# You can also use go build command for start task
# start:
#   go build -o /bin/my-app . && \
#   /bin/my-app & echo $$! > $(PID_FILE)

# Stop task will kill process by ID stored in PID_FILE (and all child processes by pstree).
stopapp:
	-kill `pstree -p \`cat $(PID_FILE)\` | tr "\n" " " |sed "s/[^0-9]/ /g" |sed "s/\s\s*/ /g"`

# Before task will only prints message. Actually, it is not necessary. You can remove it, if you want.
before:
	@echo "STOPED $(APP)" && printf '%*s\n' "40" '' | tr ' ' -

# Restart task will execute stop, before and start tasks in strict order and prints message.
restartapp: stopapp before startapp
	@echo "STARTED my-app" && printf '%*s\n' "40" '' | tr ' ' -

# Serve task will run fswatch monitor and performs restart task if any source file changed. Before serving it will execute start task.
serve: startapp
	fswatch -or --event=Updated . | \
	xargs -n1 -I {} make restartapp

# .PHONY is used for reserving tasks words
.PHONY: start before stop restart serve
