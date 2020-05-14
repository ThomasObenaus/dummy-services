.DEFAULT_GOAL := all
docker_hub_namespace := thobe
name        := ping_service

################################################################################################################
# NOTE: The following lines can keep untouched. There is nothing more to configure the docker_hub_namespace and the name.  #
#################################################################################################################

# Create version tag from git tag. Indicate if there are uncommited local changes.
flag := $(shell git diff-index --quiet HEAD -- || echo "_dirty";)
git_tag := $(shell git describe --tags -- || echo "0.0.0";)
tag  := $(git_tag)$(flag)

all: test build finish
docker: clean docker.build docker.push finish

# This target (taken from: https://gist.github.com/prwhite/8168133) is an easy way to print out a usage/ help of all make targets.
# For all make targets the text after \#\# will be printed.
help: ## Prints the help
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\1\:\2/' | column -c2 -t -s :)"

test: ## Runs the unittests
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run the unit-tests"
	@go test  -v

#-----------------
#-- build
#-----------------
build: ## builds the ping-service
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Build the $(name)"
	@go build -o $(name) .

clean:
	@rm -f version

run: build ## runs the ping-service
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run ${name}"
	@./${name}

version: delim ## creates the version file
	@echo "[INFO] Creating version file: $(tag)";
	@echo "$(tag)" > version

docker.build: delim version ## builds the docker image of the ping-service
	@echo "[INFO] Building and tagging image as $(docker_hub_namespace)/$(name):$(tag) and :latest"
	@docker build -t $(docker_hub_namespace)/$(name) --build-arg VERSION=$(tag) .
	@docker tag $(docker_hub_namespace)/$(name):latest $(docker_hub_namespace)/$(name):$(tag)

docker.push: delim ## pushes the docker image of the fail-service
	@echo "[INFO] Pushing image to Docker HUB one tagged as latest and one tagged with version :$(tag)"
	@docker login
	@docker push $(docker_hub_namespace)/$(name):$(tag)
	@docker push $(docker_hub_namespace)/$(name):latest

delim:
	@echo "------------------------------------------------------------------------------------------------"

finish:
	@echo "================================================================================================"
