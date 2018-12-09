job "fail-service" {
  datacenters = ["public-services"]

  type = "service"
  update {
    stagger = "5s"
    max_parallel = 1
  }

  group "fail-service" {
    task "fail-service" {
      driver = "docker"
      config {
        image = "<aws-account-id>.dkr.ecr.us-east-1.amazonaws.com/service/fail_service:2018-12-09_16-42-52_0ec3263_dirty"
      }


      resources {
        cpu    = 100 # MHz
        memory = 256 # MB
        network {
          mbits = 10
        }
      }
    }
  }
}