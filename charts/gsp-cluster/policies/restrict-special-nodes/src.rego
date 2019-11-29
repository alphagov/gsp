package restrict_special_nodes

violation[{"msg": msg}] {
  toleration := input.review.object.spec.tolerations[_]

  toleration.effect == "NoSchedule"

  toleration.key == input.parameters.restricted_roles[_]

  input.review.object.metadata.namespace != "gsp-system"
  input.review.object.metadata.namespace != "kube-system"

  msg := "cannot tolerate ci or cluster-management roles outside gsp-system/kube-system namespaces"
}

violation[{"msg": msg}] {
  toleration := input.review.object.spec.tolerations[_]

  not toleration.effect

  toleration.key == input.parameters.restricted_roles[_]

  input.review.object.metadata.namespace != "gsp-system"
  input.review.object.metadata.namespace != "kube-system"

  msg := "cannot tolerate ci or cluster-management roles without effect outside gsp-system/kube-system namespaces"
}

violation[{"msg": msg}] {
  toleration := input.review.object.spec.tolerations[_]

  toleration.effect == "NoSchedule"

  not toleration.key

  input.review.object.metadata.namespace != "gsp-system"
  input.review.object.metadata.namespace != "kube-system"

  msg := "cannot tolerate NoSchedule without key outside gsp-system/kube-system namespaces"
}

violation[{"msg": msg}] {
  toleration := input.review.object.spec.tolerations[_]

  not toleration.effect

  not toleration.key

  input.review.object.metadata.namespace != "gsp-system"
  input.review.object.metadata.namespace != "kube-system"

  msg := "cannot without key or effect outside gsp-system/kube-system namespaces"
}
