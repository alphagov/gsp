#!/usr/bin/env bash

set -eu

script_dir=$(dirname $0)

function usage() {
	echo "Usage: ${0} [create|destroy|template]"
}

if [ $# -lt 1 ]; then
	usage
	exit 1
fi

function log() {
	echo "💻  ${1}" 1>&2
}

function template() {
	helm template \
		"charts/${1}/" \
		--name=gsp \
		--namespace="${2}"\
		--output-dir="${3}" \
		--values="${script_dir}/local-values.yaml"
}

function template_all() {
	template gsp-cluster gsp-system "${1}"
	template gsp-istio istio-system "${1}"
	# Because we don't do Helm properly the special helm testing annotations
	# don't work. This means the test resources aren't applied at the end of an
	# install but rather immediately. This causes the test pods to error
	# because the chart won't have finished installing by the time the test
	# runs.
	rm -rf "${1}/gsp-cluster/charts/prometheus-operator/charts/grafana/templates/tests/"
}

option=${1}
cluster_name=gsp-local

case ${option} in
	destroy|delete)
		log "Destroying local GSP..."
		minikube delete --profile ${cluster_name}
		exit 0
		;;
	template)
		template_dir=${2:-manifests}
		mkdir -p "${template_dir}"
		template_all "${template_dir}"
		exit 0
		;;
	create)
		;;
	*)
		usage
		log "Unrecognised option '${option}'."
		exit 1
		;;
esac

function apply() {
	local sleep_for=10

	local apply_attempt=1
	log "[Apply attempt #${apply_attempt}] Applying ${1}..."
	until kubectl apply -R -f "${1}"; do
		log "[Apply attempt #${apply_attempt}] Failed to apply ${1}. Retrying in ${sleep_for}s..."
		apply_attempt=$((apply_attempt+1))
		sleep ${sleep_for}
	done
	log "[Apply attempt #${apply_attempt}] Successfully applied ${1}."

	sleep ${sleep_for}

	local stabilize_attempt=1
	log "[Stabilize attempt #${stabilize_attempt}] Waiting for ${1} to stabilize..."
	until ! kubectl get pods --all-namespaces | tail -n "+2" | grep -v Completed | grep -v Running; do
		log "[Stabilize attempt #${stabilize_attempt}] Failed to stabilize. Retrying in ${sleep_for}s..."
		stabilize_attempt=$((stabilize_attempt+1))
		sleep ${sleep_for}
	done

	log "[Apply attempt #${apply_attempt}, Stabilize attempt: #${stabilize_attempt}] Finished deploying ${1}."
}

function disable_minikube_dashboard() {
	if minikube addons list | grep dashboard | grep -q enabled; then
		log "[HACK] Disabling conflicting, built-in, minikube Kubernetes dashboard..."
		minikube start
		minikube addons disable dashboard
		minikube stop
		minikube delete --profile minikube
	fi
}

disable_minikube_dashboard

# Linux doesn't have hyperkit
VM_DRIVER="virtualbox"

log "Creating local GSP..."
minikube start \
	--memory ${GSP_MEMORY:-8192} \
	--cpus ${GSP_CPUS:-4} \
	--disk-size 30g \
	--vm-driver ${VM_DRIVER:-hyperkit} \
	--kubernetes-version v1.12.0 \
	--insecure-registry "registry.local.govsandbox.uk:80" \
	--profile ${cluster_name}

kubectl config set-context --current --namespace gsp-system

manifest_dir=$(mktemp -d)
#manifest_dir="/home/steve/tmp"
function cleanup() {
	rm -rf "${manifest_dir}"
	exit 0
}

# Breaks in Linux :shrug:
#trap 'cleanup' INT TERM EXIT

template_all "${manifest_dir}"

log "Applying local GSP configuration..."

log "[HACK] Creating missing namespaces..."
apply "${script_dir}/hack/create-gsp-main-namespace.yaml"

log "[HACK] Applying local DNS hack..."
apply "${script_dir}/hack/make-coredns-resolve-local-to-istio-gateway.yaml"

apply "${manifest_dir}/gsp-cluster/templates/00-aws-auth/"
apply "${manifest_dir}/gsp-istio/"
apply "${manifest_dir}/gsp-cluster/"

log "[HACK] Creating Prometheus VirtualService..."
apply "${script_dir}/hack/expose-prometheus.yaml"

log "[HACK] Creating Grafana VirtualService..."
apply "${script_dir}/hack/expose-grafana.yaml"

kubectl cluster-info
log "Local GSP ready."
