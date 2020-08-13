package restrict_special_nodes

test_allow_cluster_management_gsp_system {
  input := {
    "parameters": {
      "restricted_roles": [
        "node.kubernetes.io/cluster-management"
      ]
    },
    "review": {
      "object": {
        "metadata": {
          "namespace": "gsp-system"
        },
        "spec": {
          "tolerations": [
            {
              "effect": "NoSchedule",
              "operator": "Exists",
              "key": "node.kubernetes.io/cluster-management"
            }
          ]
        }
      }
    }
  }
  results := data.restrict_special_nodes.violation with input as input
  count(results) == 0
}

test_allow_cluster_management_kube_system {
  input := {
    "parameters": {
      "restricted_roles": [
        "node.kubernetes.io/cluster-management"
      ]
    },
    "review": {
      "object": {
        "metadata": {
          "namespace": "kube-system"
        },
        "spec": {
          "tolerations": [
            {
              "effect": "NoSchedule",
              "operator": "Exists",
              "key": "node.kubernetes.io/cluster-management"
            }
          ]
        }
      }
    }
  }
  results := data.restrict_special_nodes.violation with input as input
  count(results) == 0
}

test_deny_cluster_management_exists_main {
  input := {
    "parameters": {
      "restricted_roles": [
        "node.kubernetes.io/cluster-management"
      ]
    },
    "review": {
      "object": {
        "metadata": {
          "namespace": "sandbox-main"
        },
        "spec": {
          "tolerations": [
            {
              "effect": "NoSchedule",
              "operator": "Exists",
              "key": "node.kubernetes.io/cluster-management"
            }
          ]
        }
      }
    }
  }
  results := data.restrict_special_nodes.violation with input as input
  count(results) == 1
}

test_deny_cluster_management_equal_main {
  input := {
    "parameters": {
      "restricted_roles": [
        "node.kubernetes.io/cluster-management"
      ]
    },
    "review": {
      "object": {
        "metadata": {
          "namespace": "sandbox-main"
        },
        "spec": {
          "tolerations": [
            {
              "effect": "NoSchedule",
              "operator": "Equal",
              "value": "",
              "key": "node.kubernetes.io/cluster-management"
            }
          ]
        }
      }
    }
  }
  results := data.restrict_special_nodes.violation with input as input
  count(results) == 1
}

test_deny_cluster_management_no_effect_main {
  input := {
    "parameters": {
      "restricted_roles": [
        "node.kubernetes.io/cluster-management"
      ]
    },
    "review": {
      "object": {
        "metadata": {
          "namespace": "sandbox-main"
        },
        "spec": {
          "tolerations": [
            {
              "key": "node.kubernetes.io/cluster-management"
            }
          ]
        }
      }
    }
  }
  results := data.restrict_special_nodes.violation with input as input
  count(results) == 1
}

test_deny_no_key_main {
  input := {
    "parameters": {
      "restricted_roles": [
        "node.kubernetes.io/cluster-management"
      ]
    },
    "review": {
      "object": {
        "metadata": {
          "namespace": "sandbox-main"
        },
        "spec": {
          "tolerations": [
            {
              "effect": "NoSchedule",
            }
          ]
        }
      }
    }
  }
  results := data.restrict_special_nodes.violation with input as input
  count(results) == 1
}
