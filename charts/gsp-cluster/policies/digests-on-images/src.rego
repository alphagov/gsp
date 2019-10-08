package digests_on_images

violation[{"msg": msg}] {
  image := input.review.object.spec.containers[_].image
  registry := input.parameters.registry

  startswith(image, registry)
  not re_match("^.*@sha256:[a-f,0-9]{64}$", image)

  msg := sprintf("images from harbor must use digest (https://github.com/alphagov/gsp/blob/master/docs/gds-supported-platform/internal-images-require-digests.md): %v", [image])
}
