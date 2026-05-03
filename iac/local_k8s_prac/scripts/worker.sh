#!/usr/bin/env bash
set -euo pipefail

JOIN_FILE="/vagrant/configs/join.sh"
SENTINEL="/etc/kubernetes/kubelet.conf"

if [ -f "${SENTINEL}" ]; then
  exit 0
fi

for i in $(seq 1 60); do
  if [ -f "${JOIN_FILE}" ]; then
    break
  fi
  sleep 10
done

if [ ! -f "${JOIN_FILE}" ]; then
  echo "join file not found: ${JOIN_FILE}" >&2
  exit 1
fi

bash "${JOIN_FILE}"