package isolate_tenant_istio_resources

test_allow_correct_specification {
  input := {
    "review": {
      "object": {
        "metadata": {
          "namespace": "sometenant"
        },
        "spec": {
          "exportTo": ["."]
        }
      }
    }
  }
  results := data.isolate_tenant_istio_resources.violation with input as input
  count(results) == 0
}

test_deny_exportto_single_value {
  input := {
    "review": {
      "object": {
        "metadata": {
          "namespace": "sometenant"
        },
        "spec": {
          "exportTo": "."
        }
      }
    }
  }
  results := data.isolate_tenant_istio_resources.violation with input as input
  count(results) >= 1
}

test_deny_multiple_exportto_values {
  input := {
    "review": {
      "object": {
        "metadata": {
          "namespace": "sometenant"
        },
        "spec": {
          "exportTo": [
              ".",
              "*"
          ]
        }
      }
    }
  }
  results := data.isolate_tenant_istio_resources.violation with input as input
  count(results) >= 1
}

test_deny_exportto_star {
  input := {
    "review": {
      "object": {
        "metadata": {
          "namespace": "sometenant"
        },
        "spec": {
          "exportTo": ["*"]
        }
      }
    }
  }
  results := data.isolate_tenant_istio_resources.violation with input as input
  count(results) >= 1
}

test_deny_exportto_random_string {
  input := {
    "review": {
      "object": {
        "metadata": {
          "namespace": "sometenant"
        },
        "spec": {
          "exportTo": ["slfkjefgl"]
        }
      }
    }
  }
  results := data.isolate_tenant_istio_resources.violation with input as input
  count(results) >= 1
}

test_deny_exportto_unset {
  input := {
    "review": {
      "object": {
        "metadata": {
          "namespace": "sometenant"
        },
        "spec": {
        }
      }
    }
  }
  results := data.isolate_tenant_istio_resources.violation with input as input
  count(results) >= 1
}
