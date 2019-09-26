package digests_on_images

test_allow_docker_hub_images {
  input := {
    "parameters": {
      "registry": "registry.example.com"
    },
    "review": {
      "object": {
        "spec": {
          "containers": [{
            "image": "nginx"
          }, {
            "image": "nginx"
          }]
        }
      }
    }
  }
  results := data.digests_on_images.violation with input as input
  count(results) == 0
}

test_deny_internal_registry_with_tag {
  input := {
    "parameters": {
      "registry": "registry.example.com"
    },
    "review": {
      "object": {
        "spec": {
          "containers": [{
            "image": "nginx"
          }, {
            "image": "registry.example.com:latest"
          }]
        }
      }
    }
  }
  results := data.digests_on_images.violation with input as input
  count(results) == 1
}

test_allow_internal_registry_with_digest {
  input := {
    "parameters": {
      "registry": "registry.example.com"
    },
    "review": {
      "object": {
        "spec": {
          "containers": [{
            "image": "nginx"
          }, {
            "image": "registry.example.com@sha256:01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b"
          }]
        }
      }
    }
  }
  results := data.digests_on_images.violation with input as input
  count(results) == 0
}
