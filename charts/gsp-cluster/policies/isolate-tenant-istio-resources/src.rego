package isolate_tenant_istio_resources

violation[{"msg": msg}] {
  not input.review.object.spec.exportTo
  msg := "exportTo should be present"
}

violation[{"msg": msg}] {
  not is_array(input.review.object.spec.exportTo)
  msg := "exportTo should be a list"
}

violation[{"msg": msg}] {
  exportToCount := count(input.review.object.spec.exportTo)
  exportToCount != 1
  msg := sprintf("exportTo should be a list of size 1: %v", [exportToCount])
}

violation[{"msg": msg}] {
  exportToValue := input.review.object.spec.exportTo[0]
  exportToValue != "."
  msg := sprintf("exportTo should be set to '.': '%v'", [exportToValue])
}
