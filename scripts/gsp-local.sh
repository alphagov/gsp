#!/usr/bin/env bash

set -eu

if [ $# -lt 1 ]; then
	echo "Usage $0 [create|destroy|template]"
	exit 1
fi

function log() {
	echo "☁️  ${1}" 1>&2
}

function template() {
	helm template \
		"charts/${1}/" \
		--name=gsp \
		--namespace="${2}"\
		--output-dir="${3}" \
		--values="scripts/local-values.yaml"
}

function template_all() {
	template gsp-cluster gsp-system "${1}"
	template gsp-istio istio-system "${1}"
}

OPTION=${1}
CLUSTER_NAME=gsp-local

case ${OPTION} in
	destroy|delete)
		log "Destroying local GSP..."
		kind delete cluster --name ${CLUSTER_NAME}
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
		log "Unrecognised option '${OPTION}'."
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

	local stabilize_attempt=1
	log "[Stabilize attempt #${stabilize_attempt}] Waiting for ${1} to stabilize..."
	until ! kubectl get pods --all-namespaces | tail -n "+2" | grep -v Completed | grep -v Running; do
		log "[Stabilize attempt #${stabilize_attempt}] Failed to stabilize. Retrying in ${sleep_for}s..."
		stabilize_attempt=$((stabilize_attempt+1))
		sleep ${sleep_for}
	done

	log "[Apply attempt #${apply_attempt}, Stabilize attempt: #${stabilize_attempt}] Finished deploying ${1}."
}

log "Creating local GSP..."

kind create cluster \
	--name ${CLUSTER_NAME} \
	--image kindest/node:v1.12.5 \
	|| (log "Local GSP cluster already exists." && exit 1)

export KUBECONFIG="$(kind get kubeconfig-path --name="${CLUSTER_NAME}")"
kubectl config set-context --current  --namespace gsp-system

MANIFEST_DIR=$(mktemp -d)
function cleanup() {
	rm -rf "${MANIFEST_DIR}"
	exit 0
}
trap 'cleanup' INT TERM EXIT
template_all "${MANIFEST_DIR}"

log "Applying local GSP configuration..."

# HACK HACK HACK
kubectl apply -R -f <(cat <<EOF
apiVersion: v1
kind: Namespace
metadata:
    name: gsp-main
EOF
)

apply "${MANIFEST_DIR}/gsp-cluster/templates/00-aws-auth/"
apply "${MANIFEST_DIR}/gsp-istio/"
apply "${MANIFEST_DIR}/gsp-cluster/"

kubectl cluster-info
log "Local GSP ready."
echo "export KUBECONFIG=\"\$(kind get kubeconfig-path --name=\"${CLUSTER_NAME}\")\""
