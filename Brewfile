tap "homebrew/core"
tap "homebrew/cask"

brew "kubernetes-cli"
brew "kubernetes-helm"
brew "linuxbrew/extra/minikube" if OS.linux?
brew "opa"

if OS.mac?
  brew "hyperkit"
  brew "docker-machine-driver-hyperkit"
end

cask "minikube"
