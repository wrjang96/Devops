#!/usr/bin/env bash
set -euo pipefail

CONTROL_PLANE_IP="${CONTROL_PLANE_IP:-192.168.56.10}"
POD_CIDR="${POD_CIDR:-192.168.0.0/16}"
KUBECONFIG_DIR="/home/vagrant/.kube"
KUBECONFIG_FILE="${KUBECONFIG_DIR}/config"
SENTINEL="/etc/kubernetes/admin.conf"
JOIN_FILE="/vagrant/configs/join.sh"
HOST_KUBECONFIG="/vagrant/configs/config"

if [ ! -f "${SENTINEL}" ]; then
  sudo kubeadm init \
    --apiserver-advertise-address="${CONTROL_PLANE_IP}" \
    --pod-network-cidr="${POD_CIDR}"
fi

mkdir -p "${KUBECONFIG_DIR}"
sudo cp -f /etc/kubernetes/admin.conf "${KUBECONFIG_FILE}"
sudo chown vagrant:vagrant "${KUBECONFIG_FILE}"

export KUBECONFIG="${KUBECONFIG_FILE}"

if ! kubectl get daemonset -n kube-system calico-node >/dev/null 2>&1; then
  kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.30.3/manifests/calico.yaml
fi

sudo kubeadm token create --print-join-command > /tmp/join.sh
echo "sudo $(cat /tmp/join.sh) --ignore-preflight-errors=Swap" | tee "${JOIN_FILE}" >/dev/null
chmod +x "${JOIN_FILE}"

cp -f "${KUBECONFIG_FILE}" "${HOST_KUBECONFIG}"
chmod 644 "${HOST_KUBECONFIG}"