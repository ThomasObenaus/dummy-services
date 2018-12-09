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
        image = "<aws-account-id>.dkr.ecr.us-east-1.amazonaws.com/service/fail_service:2018-12-09_22-14-45_22f3e74_dirty"
        port_map = {
          http = 8080
        }
      }

      # Register at consul
      service {
        name = "${TASK}"
        port = "http"
        check {
          port     = "http"
          type     = "http"
          path     = "/health"
          method   = "GET"
          interval = "10s"
          timeout  = "2s"
        }
      }

      env {
        HEALTHY_IN    = 20,
        HEALTHY_FOR   = 50,
        UNHEALTHY_FOR = 100,
      }

      resources {
        cpu    = 100 # MHz
        memory = 256 # MB
        network {
          mbits = 10
          port "http" {}
        }
      }
    }
  }
}