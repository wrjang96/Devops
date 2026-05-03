# CKAD Local Kubernetes Lab Integrated Guide

## 1) 목표 구성

- macOS Apple Silicon
- `Vagrant + VirtualBox`
- Ubuntu 24.04 ARM64 VM
- `containerd + kubeadm + kubelet + kubectl`
- CNI: Calico
- 기본 노드: `control-plane` 1대, `worker1` 1대

## 2) 사전 준비

필수 도구 확인:

```bash
brew --version
vagrant --version
VBoxManage --version
kubectl version --client
git --version
uname -m
```

필요 시 설치:

```bash
brew install --cask virtualbox
brew install --cask vagrant
brew install kubectl git
```

Apple Silicon이면 `uname -m` 결과가 `arm64`입니다.

현재 기본 box:

```bash
export BOX_NAME="bento/ubuntu-24.04"
```

## 3) 현재 파일 구조

현재 구조는 `lab/` 하위가 아니라 `iac/local_k8s_prac/` 바로 아래에 있습니다.

```text
iac/local_k8s_prac/
├── Makefile
├── README.md
├── Vagrantfile
├── configs/
│   ├── .gitkeep
│   ├── config       # generated, ignored
│   └── join.sh      # generated, ignored
├── exercises/
│   └── README.md
└── scripts/
    ├── common.sh
    ├── control-plane.sh
    ├── host-export.sh
    └── worker.sh
```

## 4) Vagrantfile 역할

`Vagrantfile`은 VM을 정의합니다.

- `BOX_NAME`: 기본 `bento/ubuntu-24.04`
- `PROVIDER`: 기본 `virtualbox`
- `control-plane`: `192.168.56.10`
- `worker1`: `192.168.56.11`
- `ENABLE_SECOND_WORKER=true`이면 `worker2` 추가

문법 확인:

```bash
ruby -c Vagrantfile
```

기대 출력:

```text
Syntax OK
```

## 5) Provisioning 스크립트 역할

`scripts/common.sh`

- 모든 VM에서 공통 실행
- swap 비활성화
- `overlay`, `br_netfilter` 모듈 및 sysctl 설정
- `containerd` 설치 및 `SystemdCgroup = true` 설정
- Kubernetes apt repo 등록
- `kubelet`, `kubeadm`, `kubectl` 설치
- VM 내부에서 `k=kubectl` alias 등록

`scripts/control-plane.sh`

- `kubeadm init` 실행
- `/home/vagrant/.kube/config` 생성
- Calico 설치
- worker join 명령을 `/vagrant/configs/join.sh`에 저장
- host Mac용 kubeconfig를 `/vagrant/configs/config`에 복사

`scripts/worker.sh`

- `/vagrant/configs/join.sh`가 생길 때까지 대기
- join 파일을 실행해 worker를 클러스터에 붙임

`scripts/host-export.sh`

- host Mac에서 `configs/config`를 `KUBECONFIG`로 지정
- `kubectl cluster-info`
- `kubectl get nodes -o wide`

참고:

- `host-export.sh` 안의 `export KUBECONFIG=...`는 스크립트 안에서만 유지됩니다.
- 이후 로컬 터미널에서 계속 `kubectl`을 쓰려면 직접 `export KUBECONFIG=$(pwd)/configs/config`를 실행해야 합니다.

## 6) 클러스터 생성

작업 디렉터리 이동:

box가 이미 있는지 확인:

```bash
vagrant box list
```

VM 생성 및 provisioning:

```bash
vagrant up
```

특정 노드만 올리기:

```bash
vagrant up control-plane
vagrant up worker1
```

상태 확인:

```bash
vagrant status
```

## 7) 클러스터 확인

control-plane VM 안에서:

```bash
vagrant ssh control-plane
kubectl get nodes
kubectl get pods -A
k get nodes
exit


host Mac에서 계속 `kubectl`을 쓰려면:

```bash
export KUBECONFIG={local directory}/configs/config
kubectl get nodes -o wide
kubectl get pods -A
```

정상 기준:

```text
control-plane   Ready
worker1         Ready
```

주의:

- VirtualBox NAT 때문에 `kubectl get nodes -o wide`에서 두 노드의 `INTERNAL-IP`가 `10.0.2.15`로 보일 수 있습니다.
- 기본 실습은 가능하지만, 더 정확히 하려면 kubelet node IP를 private IP로 지정하는 개선이 필요합니다.

## 8) CKAD 기본 검증

namespace와 Pod 생성/삭제:

```bash
kubectl create namespace ckad-practice
kubectl get ns
kubectl run nginx --image=nginx -n ckad-practice
kubectl get pods -n ckad-practice
kubectl delete namespace ckad-practice
```

간단한 Deployment/Service:

```bash
kubectl create deployment web --image=nginx
kubectl expose deployment web --port=80 --type=ClusterIP
kubectl get deploy,svc,pods
kubectl delete deployment web
kubectl delete service web
```

## 9) Makefile 단축 명령

`Makefile`은 프로젝트 전용 alias 모음처럼 쓰입니다.

```bash
make status
make nodes
make ssh-cp
make ssh-worker1
make halt
make up
make destroy
```

의미:

- `make status`: `vagrant status`
- `make nodes`: `configs/config`로 `kubectl get nodes -o wide`
- `make ssh-cp`: `vagrant ssh control-plane`
- `make ssh-worker1`: `vagrant ssh worker1`
- `make halt`: VM 전원 끄기
- `make up`: VM 생성/기동
- `make destroy`: VM 삭제

주의:

- `make destroy`는 VM 디스크를 삭제합니다.
- VM 안에서 직접 수정한 설정은 사라집니다.
- 스크립트에 넣어둔 alias나 provisioning 설정은 다시 `vagrant up` 때 재적용됩니다.

## 10) 운영 점검

host Mac에서:

```bash
export KUBECONFIG={local directory}/configs/config
kubectl get nodes -o wide
kubectl get pods -A
kubectl cluster-info
```

control-plane VM에서:

```bash
vagrant ssh control-plane
sudo systemctl status containerd --no-pager
sudo systemctl status kubelet --no-pager
kubectl get pods -A
exit
```

worker VM에서:

```bash
vagrant ssh worker1
sudo systemctl status containerd --no-pager
sudo systemctl status kubelet --no-pager
exit
```

`--no-pager`는 `less` 같은 pager 화면으로 들어가지 않고 결과를 터미널에 바로 출력하라는 옵션입니다.

## 11) Calico 확인

```bash
export KUBECONFIG={local directory}/configs/config
kubectl get pods -n kube-system | grep calico
```

또는 전체 kube-system 확인:

```bash
kubectl get pods -n kube-system
```

Calico는 host Mac에 설치되는 것이 아니라 VM 안 Kubernetes 클러스터에 CNI로 설치됩니다.

## 12) 재기동과 재생성

VM 전원 끄기:

```bash
vagrant halt
```

다시 올리기:

```bash
vagrant up
```

재프로비저닝:

```bash
vagrant provision
vagrant provision control-plane
vagrant provision worker1
```

worker만 다시 만들기:

```bash
vagrant destroy -f worker1
vagrant up worker1
```

전체 초기화:

```bash
vagrant destroy -f
rm -f configs/config configs/join.sh
vagrant up
```

## 13) kubeconfig 복구

`configs/config`가 없으면 control-plane에서 다시 복사합니다.

```bash
vagrant ssh control-plane
mkdir -p /vagrant/configs
cp ~/.kube/config /vagrant/configs/config
chmod 644 /vagrant/configs/config
exit
```

확인:

```bash
ls -l configs
./scripts/host-export.sh
```