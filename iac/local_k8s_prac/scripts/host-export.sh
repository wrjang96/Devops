#!/usr/bin/env bash
set -euo pipefail

KUBECONFIG_PATH="${PWD}/configs/config"

if [ ! -f "${KUBECONFIG_PATH}" ]; then
  echo "missing kubeconfig: ${KUBECONFIG_PATH}" >&2
  exit 1
fi

export KUBECONFIG="${KUBECONFIG_PATH}"
kubectl cluster-info
kubectl get nodes -o wide