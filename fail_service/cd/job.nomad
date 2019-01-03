job "fail-service" {
  datacenters = ["testing"]

  type = "service"

  update {
    stagger = "5s"
    max_parallel = 1
  }

  group "fail-service" {
    task "fail-service" {
      driver = "docker"
      config {
        image = "thobe/fail_service:latest"
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
        HEALTHY_IN    = 0,
        HEALTHY_FOR   = 0,
        UNHEALTHY_FOR = 0,
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