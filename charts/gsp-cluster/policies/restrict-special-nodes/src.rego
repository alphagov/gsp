package restrict_special_nodes

violation[{"msg": msg}] {
  toleration := input.review.object.spec.tolerations[_]

  toleration.effect == "NoSchedule"
  toleration.operator == "Exists"

  toleration.key == input.parameters.restricted_roles[_]

  input.review.object.metadata.namespace != "gsp-system"
  input.review.object.metadata.namespace != "kube-system"

  msg := "cannot tolerate ci or cluster-management roles outside gsp-system/kube-system namespaces"
}
