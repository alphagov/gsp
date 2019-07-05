tap "homebrew/core"
tap "homebrew/cask"

brew "kubernetes-cli"
brew "kubernetes-helm"

if OS.mac?
  brew "hyperkit"
  brew "docker-machine-driver-hyperkit"
  cask "aws-vault"
  cask "minikube"
else
  brew "linuxbrew/extra/aws-vault"
  brew "linuxbrew/extra/minikube"
end
