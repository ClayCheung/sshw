# sshw

[![Build Status](https://travis-ci.org/yinheli/sshw.svg?branch=master)](https://travis-ci.org/yinheli/sshw)

ssh client wrapper for automatic login.

![usage](./assets/sshw-demo.gif)

## install

use `go get`

```
go get -u github.com/yinheli/sshw/cmd/sshw
```

or download binary from [releases](//github.com/yinheli/sshw/releases).

## config

put config file in `~/.sshw` or `~/.sshw.yml` or `~/.sshw.yaml` or `./.sshw` or `./.sshw.yml` or `./.sshw.yaml`.

config example:

```yaml
- { name: dev server fully configured, user: appuser, host: 192.168.8.35, port: 22, password: 123456 }
- { name: dev server with key path, user: appuser, host: 192.168.8.35, port: 22, keypath: /root/.ssh/id_rsa }
- { name: dev server with passphrase key, user: appuser, host: 192.168.8.35, port: 22, keypath: /root/.ssh/id_rsa, passphrase: abcdefghijklmn}
- { name: dev server without port, user: appuser, host: 192.168.8.35 }
- { name: dev server without user, host: 192.168.8.35 }
- { name: dev server without password, host: 192.168.8.35 }
- { name: ⚡️ server with emoji name, host: 192.168.8.35 }
- { name: server with alias, alias: dev, host: 192.168.8.35 }
- name: server with jump
  user: appuser
  host: 192.168.8.35
  port: 22
  password: 123456
  jump:
  - user: appuser
    host: 192.168.8.36
    port: 2222


# server group 1
- name: server group 1
  children:
  - { name: server 1, user: root, host: 192.168.1.2 }
  - { name: server 2, user: root, host: 192.168.1.3 }
  - { name: server 3, user: root, host: 192.168.1.4 }

# server group 2
- name: server group 2
  children:
  - { name: server 1, user: root, host: 192.168.2.2 }
  - { name: server 2, user: root, host: 192.168.3.3 }
  - { name: server 3, user: root, host: 192.168.4.4 }
  
# For K8S Cluster
# first feature: if host set like 192.168.1.2(dev.harbor.io), it will login by 192.168.1.2,
# and set `192.168.1.2 dev.harbor.io` in `/etc/hosts` (after you use `sudo sshw` to login)
#
# second feature: if set `kube: xxx`, it will fetch kubeConfig on remote server, 
# and set it on local kubeConfig, so you can use `kubectl --context=xxx` to access remote k8s cluster
- name: k8s-dev
  children:
    - { name: registry, user: root, host: 192.168.1.2(dev.harbor.io), password: 123456 }
    - { name: master-1, user: root, host: 192.168.1.3, password: 123456, kube: k8s-dev }
    - { name: master-2, user: root, host: 192.168.1.4, password: 123456 }
    - { name: master-3, user: root, host: 192.168.1.5, password: 123456 }
    - { name: node-1, user: root, host: 192.168.1.6, password: 123456 }
    - { name: node-2, user: root, host: 192.168.1.7, password: 123456 }
```

# callback
```
- name: dev server fully configured
  user: appuser
  host: 192.168.8.35
  port: 22
  password: 123456
  callback-shells:
  - {cmd: 2}
  - {delay: 1500, cmd: 0}
  - {cmd: 'echo 1'}
 ```
