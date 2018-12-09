.DEFAULT_GOAL := all
docker_hub_namespace := thobe
name        := fail_service
category    := service
aws_reg     := us-east-1
aws_profile := playground

################################################################################################################
# NOTE: The following lines can keep untouched. There is nothing more to configure the docker_hub_namespace and the name.  #
#################################################################################################################

# obtain aws account id
aws_aid     := $(shell aws sts get-caller-identity --output text --query 'Account' --profile $(aws_profile))
ecr_url  := $(aws_aid).dkr.ecr.$(aws_reg).amazonaws.com

# Create credentials for Docker for AWS ecr login
creds := $(shell aws ecr get-login --no-include-email --region $(aws_reg) --profile $(aws_profile))

# Create version tag from git commit message. Indicate if there are uncommited local changes.
date := $(shell date '+%Y-%m-%d_%H-%M-%S')
rev  := $(shell git rev-parse --short HEAD)
flag := $(shell git diff-index --quiet HEAD -- || echo "_dirty";)
tag  := $(date)_$(rev)$(flag)


all: vendor test build finish
docker.aws: clean version docker.aws.build docker.aws.push finish
docker.hub: clean version docker.hub.build docker.hub.push finish

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
	@./${name} --healthy-for 20 --healthy-in 10 --unhealthy-for 10

version: delim
	@echo "[INFO] Building version:"
	@echo "$(tag)" | tee version

docker.hub.build: delim
	@echo "[INFO] Building and tagging image as $(docker_hub_namespace)/$(name):$(tag)"
	@docker build -t $(docker_hub_namespace)/$(name) --build-arg VERSION=$(tag) .
	@docker tag $(docker_hub_namespace)/$(name):latest $(docker_hub_namespace)/$(name):$(tag)

docker.hub.push: delim
	@echo "[INFO] Pushing image to Docker HUB"
	@docker login
	@docker push $(docker_hub_namespace)/$(name):$(tag)

credentials: delim
	@echo "[INFO] Login to AWS ECR"
	@$(creds)

docker.aws.build: delim
	@echo "[INFO] Building and tagging image as $(ecr_url)/$(category)/$(name):$(tag)"
	@docker build -t $(category)/$(name) --build-arg VERSION=$(tag) .
	@docker tag $(category)/$(name):latest $(ecr_url)/$(category)/$(name):$(tag)

docker.aws.push: credentials delim
	@echo "[INFO] Pushing image to AWS ECR"
	@docker push $(ecr_url)/$(category)/$(name):$(tag)

delim:
	@echo "------------------------------------------------------------------------------------------------"

finish:
	@echo "================================================================================================"
