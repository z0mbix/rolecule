# Rolecule

## Description

`rolecule` is a simple tool to help you test your configuration management code
works as you expect, but creating systemd enabled containers with either docker or podman, then converging them with your configured provisioner (ansible by default).

Once converged, it will run a verifier to test it all. Currently, then only supported provisioner is [goss](https://github.com/goss-org/goss), [testinfra](https://testinfra.readthedocs.io/) will be added soon.

## Usage

First, you need to create a simple `rolecule.yml` file in the root of your role/module/recipe, e.g.:

```yaml
---
engine:
  name: podman

containers:
  - name: rockylinux-systemd-9.1-amd64
    image: localhost/rockylinux-systemd:9.1
    arch: amd64

provisioner:
  name: ansible
  command: ansible-playbook --connection local --inventory localhost,
  env:
    ANSIBLE_ROLES_PATH: .
    ANSIBLE_NOCOWS: True

verifier:
  name: goss
```

Then, from the root of the role, run `rolecule test`, e.g.:

```shell

```
