.DEFAULT_GOAL := all
docker_hub_namespace := thobe
name        := fail_service

################################################################################################################
# NOTE: The following lines can keep untouched. There is nothing more to configure the docker_hub_namespace and the name.  #
#################################################################################################################

# Create version tag from git commit message. Indicate if there are uncommited local changes.
date := $(shell date '+%Y-%m-%d_%H-%M-%S')
rev  := $(shell git rev-parse --short HEAD)
flag := $(shell git diff-index --quiet HEAD -- || echo "_dirty";)
tag  := $(date)_$(rev)$(flag)


all: vendor test build finish
docker: clean version docker.build docker.push finish

test:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run the unit-tests"
	@go test  -v

#-----------------
#-- build
#-----------------
build:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Build the $(name)"
	@go build -o $(name) .

#------------------
#-- dependencies
#------------------
vendor: depend.install depend.update

depend.update:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> updating dependencies from Gopkg.lock"
	@dep ensure -update -v

depend.install:
	@echo "----------------------------------------------------------------------------------"
	@echo "--> install dependencies as listed in Gopkg.toml"
	@dep ensure -v

clean:
	@rm -f version

run: build
	@echo "----------------------------------------------------------------------------------"
	@echo "--> Run ${name}"
	@./${name}

version: delim
	@echo "[INFO] Building version:"
	@echo "$(tag)" | tee version

docker.build: delim
	@echo "[INFO] Building and tagging image"
	@docker build -t $(docker_hub_namespace)/$(name) --build-arg VERSION=$(tag) .
	@docker tag $(docker_hub_namespace)/$(name):latest $(docker_hub_namespace)/$(name):$(tag)

docker.push: delim
	@echo "[INFO] Pushing image to Docker HUB"
	@docker login
	@docker push $(docker_hub_namespace)/$(name):$(tag)

delim:
	@echo "------------------------------------------------------------------------------------------------"

finish:
	@echo "================================================================================================"
