#!/usr/bin/env bash
set -euo pipefail

K8S_VERSION="${K8S_VERSION:-1.35}"

sudo swapoff -a
sudo sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab

cat <<SYSCTL | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
SYSCTL

sudo modprobe overlay
sudo modprobe br_netfilter

cat <<SYSCTL | sudo tee /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward = 1
SYSCTL

sudo sysctl --system

sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gpg

if ! command -v containerd >/dev/null 2>&1; then
  sudo apt-get install -y containerd
fi

sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml >/dev/null
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
sudo systemctl restart containerd
sudo systemctl enable containerd

sudo mkdir -p /etc/apt/keyrings
curl -fsSL "https://pkgs.k8s.io/core:/stable:/v${K8S_VERSION}/deb/Release.key" \
  | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v${K8S_VERSION}/deb/ /" \
  | sudo tee /etc/apt/sources.list.d/kubernetes.list

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
sudo systemctl enable kubelet

echo 'alias k=kubectl' >> /home/vagrant/.bashrc
echo 'source <(kubectl completion bash)' >> /home/vagrant/.bashrc
echo 'complete -o default -F __start_kubectl k' >> /home/vagrant/.bashrc
