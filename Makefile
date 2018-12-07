.DEFAULT_GOAL				:= all
category    := service
name        := ping-service
aws_reg     := eu-central-1
aws_profile := playground

################################################################################################################
# NOTE: The following lines can keep untouched. There is nothing more to configure the category and the name.  #
#################################################################################################################

# obtain aws account id
aws_aid     := $(shell aws sts get-caller-identity --output text --query 'Account' --profile $(aws_profile))

ecr_url  := $(aws_aid).dkr.ecr.$(aws_reg).amazonaws.com

# Create version tag from git commit message. Indicate if there are uncommited local changes.
date := $(shell date '+%Y-%m-%d_%H-%M-%S')
rev  := $(shell git rev-parse --short HEAD)
flag := $(shell git diff-index --quiet HEAD -- || echo "_dirty";)
tag  := $(date)_$(rev)$(flag)

# Create credentials for Docker for AWS ecr login
creds := $(shell aws ecr get-login --no-include-email --region $(aws_reg) --profile $(aws_profile))

all: vendor test build finish
all-docker: clean vendor test docker push finish

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

credentials: delim
	@echo "[INFO] Login to AWS ECR"
	@$(creds)

docker: delim
	@echo "[INFO] Building and tagging image"
	docker build -t $(category)/$(name) --build-arg VERSION=$(tag) .
	@docker tag $(category)/$(name):latest $(ecr_url)/$(category)/$(name):$(tag)

push: credentials delim
	@echo "[INFO] Pushing image to AWS ECR"
	@docker push $(ecr_url)/$(category)/$(name):$(tag)

delim:
	@echo "------------------------------------------------------------------------------------------------"

finish:
	@echo "================================================================================================"
