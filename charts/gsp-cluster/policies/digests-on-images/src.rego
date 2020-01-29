package digests_on_images

violation[{"msg": msg}] {
  image := input.review.object.spec.containers[_].image
  registry := input.parameters.registry

  startswith(image, registry)
  not re_match("^.*@sha256:[a-f,0-9]{64}$", image)

  msg := sprintf("images from harbor must use digest (https://github.com/alphagov/gsp/blob/master/docs/gds-supported-platform/internal-images-require-digests.md): %v", [image])
}

violation[{"msg": msg}] {
  image := input.review.object.spec.containers[_].image
  aws_account_id := input.parameters.aws_account_id
  registry = sprintf("%s.dkr.ecr.eu-west-2.amazonaws.com", [aws_account_id])

  startswith(image, registry)
  not re_match("^.*@sha256:[a-f,0-9]{64}$", image)

  msg := sprintf("images from ecr must use digest (https://github.com/alphagov/gsp/blob/master/docs/gds-supported-platform/internal-images-require-digests.md): %v", [image])
}
