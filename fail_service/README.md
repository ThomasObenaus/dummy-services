# Fail-Service

This is a small service for testing purposes. The only thin it can do is get healthy or unhealthy.
Therefore it provides a `/health` endpoint. It reports 200_OK if the service is healthy or 504_GatewayTimeout otherwise.

The service state can be influenced via command line parameters or by sending a request to the `/sethealthy` or `/setunhealthy` endpoint.

## Build and Run Locally

To build an run the service locally you can just use the according makefile target.
By calling `make run` the serive will be build and started with default settings.

These settings are printed to stdout:

```bash
2019/01/03 14:32:56 Cfg:
2019/01/03 14:32:56     healthyIn: 10
2019/01/03 14:32:56     healthyFor: 20
2019/01/03 14:32:56     unhealthyFor: 10
2019/01/03 14:32:56 Starts listening at 8080.
```

To just build and then run it with custom parameters you just call `make build`. This creates a binary called `fail_service`.

To start it you just call it with the correct parameters. I.e. `./fail_service -healthy-in=10`

## Build docker and push it to Docker Hub

- `make docker` will build the docker image and push it to Docker Hub

## Deploy it to Nomad

- Adjust file `cd/job.nomad` by setting the variable `datacenters = ["testing"]` to the data center of your nomad cluster.
- Deploy it via `nomad run cd/job.nomad`

## Command Line Interface

```bash
Usage of ./fail_service:
  -healthy-for=0: Number of seconds the health end-point will return a 200. 0 means forever.
  -healthy-in=0: Number of seconds the health end-point will start returning a 200.
  -p=8080: The port where the application instance listens to. Defaults to 8080.
  -unhealthy-for=0: Number of seconds the health end-point will keep returning a !200.
```

### Examples

```bash
# Starts healthy and stays healthy
./fail_service

# Gets healthy in 10s and stays healthy
./fail_service -healthy-in=10

# Starts healthy, then after 20s it gets unhealthy. For 3s it stays unhealthy
# and gets healthy again to stay so for 20s. etc.
./fail_service -healthy-for=20 -unhealthy-for=3

# Gets healthy in 10s, then after 20s it gets unhealthy. For 3s it stays unhealthy
# and gets healthy again to stay so for 20s. etc.
./fail_service -healthy-in=10 -healthy-for=20 -unhealthy-for=3

# Service will stay unhealthy forever
./fail_service -healthy-in=-1
```

## Overwrite State via HTTP Call

The service provides two http endpoints which can be used to set the health state of the service directly. Doing this the configured pattern defined by the CLI will be overwritten and not regarded any more. This means as soon as one of the endpoints is called this health state stays until the service is restarted.

Both endpoints expect a PUT call.

```bash
# Set the service healthy
curl -X PUT localhost:8080/sethealthy

# Set the service unhealthy
curl -X PUT localhost:8080/setunhealthy
```
