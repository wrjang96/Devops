# Laravel (Apache + PHP) + PostgreSQL on Vagrant, provisioned by Ansible

Vagrant VM(Ubuntu) 위에 Apache/PHP로 Laravel 앱을 배포하고 PostgreSQL을 구성하는 Ansible 플레이북 프로젝트입니다.  
Laravel 소스는 `src/` 디렉터리에 배포됩니다.

## 구성 요소
- Provisioning: Ansible (`playbook.yml`)
- VM: Vagrant + VirtualBox (`Vagrantfile`)
- Web: Apache2 + PHP 8.x
- DB: PostgreSQL
- App: Laravel (`src/`)

## 디렉터리 구조
- `playbook.yml` : 메인 플레이북
- `Vagrantfile` : VM 정의
- `hosts` : 실사용 `hosts`는 커밋하지 않음
- `group_vars/` : 그룹 변수
  - `vagrant.yml` : 공통 접속/환경 변수
  - `web.yml` : Apache/Laravel 웹 관련 변수
  - `db.yml` : 실사용 db.yml은 는 커밋하지 않음
- `files/`
  - `apache/laravel.conf.j2` : Apache vhost 템플릿
  - `php/zzmyphp.ini` : PHP 설정
  - `deploy/.env.j2` : Laravel `.env` 템플릿 (민감값 하드코딩 금지)
- `src/` : Laravel 애플리케이션 디렉터리 (배포 대상)

## 요구 사항
- macOS (호스트)
- VirtualBox, Vagrant
- Ansible (호스트에 설치)
- 인터넷 연결 (VM에서 apt/composer 다운로드)

## 빠른 시작

### 1) VM 생성
```bash
vagrant up
```

### 2) 인벤토리 준비
```bash
cp hosts.example hosts
```
`hosts` 안의 IP/키 경로를 본인 환경에 맞게 수정합니다.

### 3) (권장) DB 변수 파일 준비
```bash
cp group_vars/db.yml.example group_vars/db.yml
```
`group_vars/db.yml`의 `db_password` 등 민감값을 실제 값으로 설정합니다.  

### 4) 플레이북 실행
```bash
ansible-playbook -i hosts playbook.yml
```

태그 단위 실행 예시:
```bash
ansible-playbook -i hosts playbook.yml --tags "apache,php"
ansible-playbook -i hosts playbook.yml --tags "deploy"
```

## 동작 확인(호스트 macOS 기준)

### A. Apache 응답 확인
VM IP가 `192.168.34.23`이라면:
```bash
curl -I http://192.168.34.23
```

도메인으로 접근하려면 `/etc/hosts`에 추가:
```txt
192.168.34.23  laravel.dev
```
그 뒤 브라우저에서 `http://laravel.dev` 접속.

### B. VM 내부에서 서비스 상태 확인
```bash
vagrant ssh
sudo systemctl status apache2
sudo systemctl status postgresql
```

### C. Laravel 앱 동작 확인
```bash
vagrant ssh
cd /vagrant/src
php -v
php artisan --version
php artisan about
```

### D. DB 연결 확인(앱 관점)
```bash
vagrant ssh
cd /vagrant/src
php artisan migrate:status
```
- 정상이라면 마이그레이션 목록이 출력됩니다.
- `Could not open input file: artisan`이면 **현재 디렉터리가 Laravel 루트가 아닙니다.**
  - `cd /vagrant/src` 후 다시 실행하세요.

### E. PostgreSQL 직접 확인
```bash
vagrant ssh
sudo -u postgres psql -c "\l"
sudo -u postgres psql -c "\du"
```