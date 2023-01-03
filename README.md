# Rolecule

## Description

`rolecule` is a simple tool to help you test your configuration management code works as you expect, by creating systemd enabled containers with either docker or podman, then converging them with your configured provisioner (ansible by default). We're basically treating containers as mini VMs.

Once converged, it will run a verifier to test it all. Currently, the only supported provisioner is [goss](https://github.com/goss-org/goss), [testinfra](https://testinfra.readthedocs.io/) will be added soon.

This should speed up testing your roles if you're using virtual machines.

## Usage

First, you need to create a simple `rolecule.yml` file in the root of your role/module/recipe, e.g.:

```yaml
---
engine:
  name: podman

containers:
  - name: rockylinux-9.1
    image: rockylinux-systemd:9.1
  - name: ubuntu-22.04
    image: ubuntu-systemd:22.04

provisioner:
  name: ansible

verifier:
  name: goss
```

Then, from the root of your role (e.g. sshd), run `rolecule test`, e.g.:

```
» rolecule test
   • creating container rolecule-sshd-rockylinux-9.1 with podman
   • creating container rolecule-sshd-ubuntu-22.04 with podman
   • converging container rolecule-sshd-rockylinux-9.1 with ansible
Using /etc/ansible/ansible.cfg as config file

PLAY [test] ********************************************************************
...

PLAY RECAP *********************************************************************
localhost                  : ok=5    changed=4    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

   • converging container rolecule-sshd-ubuntu-22.04 with ansible
No config file found; using defaults

PLAY [test] ********************************************************************
...

PLAY RECAP *********************************************************************
localhost                  : ok=5    changed=4    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

   • verifying container rolecule-sshd-rockylinux-9.1 with goss
1..12
ok 1 - File: /etc/ssh/sshd_config: exists: matches expectation: [true]
ok 2 - File: /etc/ssh/sshd_config: mode: matches expectation: ["0600"]
ok 3 - File: /etc/ssh/sshd_config: owner: matches expectation: ["root"]
ok 4 - File: /etc/ssh/sshd_config: group: matches expectation: ["root"]
ok 5 - File: /etc/ssh/sshd_config: filetype: matches expectation: ["file"]
ok 6 - File: /etc/ssh/sshd_config: contains: all expectations found: [Port 22, AddressFamily any, ListenAddress 0.0.0.0, SyslogFacility AUTH, LogLevel INFO, PermitRootLogin no, PubkeyAuthentication yes, AuthorizedKeysFile	.ssh/authorized_keys, Subsystem	sftp	/usr/libexec/openssh/sftp-server]
ok 7 - User: sshd: exists: matches expectation: [true]
ok 8 - Process: sshd: running: matches expectation: [true]
ok 9 - Port: tcp:22: listening: matches expectation: [true]
ok 10 - Port: tcp:22: ip: matches expectation: [["0.0.0.0"]]
ok 11 - Service: sshd: enabled: matches expectation: [true]
ok 12 - Service: sshd: running: matches expectation: [true]
   • verifying container rolecule-sshd-ubuntu-22.04 with goss
1..12
ok 1 - File: /etc/ssh/sshd_config: exists: matches expectation: [true]
ok 2 - File: /etc/ssh/sshd_config: mode: matches expectation: ["0600"]
ok 3 - File: /etc/ssh/sshd_config: owner: matches expectation: ["root"]
ok 4 - File: /etc/ssh/sshd_config: group: matches expectation: ["root"]
ok 5 - File: /etc/ssh/sshd_config: filetype: matches expectation: ["file"]
ok 6 - File: /etc/ssh/sshd_config: contains: all expectations found: [Port 22, AddressFamily any, ListenAddress 0.0.0.0, SyslogFacility AUTH, LogLevel INFO, PermitRootLogin no, PubkeyAuthentication yes, AuthorizedKeysFile	.ssh/authorized_keys, Subsystem	sftp	/usr/libexec/openssh/sftp-server]
ok 7 - User: sshd: exists: matches expectation: [true]
ok 8 - Process: sshd: running: matches expectation: [true]
ok 9 - Port: tcp:22: listening: matches expectation: [true]
ok 10 - Port: tcp:22: ip: matches expectation: [["0.0.0.0"]]
ok 11 - Service: sshd: enabled: matches expectation: [true]
ok 12 - Service: sshd: running: matches expectation: [true]
   • destroying container rolecule-sshd-rockylinux-9.1
   • destroying container rolecule-sshd-ubuntu-22.04
   • complete
```

## Help

```text
» rolecule --help
rolecule uses docker or podman to test your
configuration management roles/recipes/modules in a systemd enabled container,
then tests them with a verifier (goss/testinfra).

Usage:
  rolecule [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  converge    Run your configuration management tool to converge the container(s)
  create      Create a new container(s) to test the role in
  destroy     Destroy everything
  help        Help about any command
  init        Initialise the project with a nice new rolecule.yml file
  list        List the running containers for this role/module/recipe
  shell       Open a shell in a container
  test        Create the container(s), converge them, test them, then clean up
  verify      Verify your containers are configured how you expect

Flags:
  -d, --debug   enable debug output
  -h, --help    help for rolecule

Use "rolecule [command] --help" for more information about a command.
```

## Individual sub command examples

Create the containers:

```text
» rolecule create
   • creating container rolecule-sshd-rockylinux-9.1 with podman
   • creating container rolecule-sshd-ubuntu-22.04 with podman
```

List all containers:

```text
» rolecule list
CONTAINER ID  IMAGE                             COMMAND         CREATED         STATUS             PORTS       NAMES
e9638f571210  localhost/rockylinux-systemd:9.1  /usr/sbin/init  21 minutes ago  Up 21 minutes ago              rolecule-sshd-rockylinux-9.1
0e07e4214fd5  localhost/ubuntu-systemd:22.04                    21 minutes ago  Up 21 minutes ago              rolecule-sshd-ubuntu-22.04
```

Open a shell:

```text
» rolecule shell -n rolecule-sshd-ubuntu-22.04
root@0e07e4214fd5:/src#
```

Destroy the containers:

```text
» rolecule destroy
   • destroying container rolecule-sshd-rockylinux-9.1
   • destroying container rolecule-sshd-ubuntu-22.04
```

### FAQ

**How do I get this working on macOS?**

You'll need to make sure you create a podman machine with your home directory mounted for volume mounts to work, e.g.:

```
» podman machine init --now --rootful -v $HOME:$HOME
```

**How do I create a suitable container image for this?**

You can use the `Containerfile`/`Dockerfile` files in the testing directory to build suitable images:

```
» podman build -t rockylinux-systemd:9.1 -f testing/rockylinux-9.1-systemd.Containerfile .
```

## TODO

- ~~Test on Mac~~
- ~~Test on Linux~~
- Support installing ansible collections
- Support testinfra verifier
- Support scenarios, making it possible to test a role with different tags
- Support using custom provisioner command/args/env vars from rolecule.yml
- Support using custom verifier command/args/env vars from rolecule.yml
- Test converging with puppet apply
- Implement `rolecule init` to generate a rolecule.yml file (use current directory structure to determine configuration management provisioner)
- ~~Implement `rolecule list` subcommand to list all running containers~~
- Write some tests :/
- Document what is required for a container image
- Test with docker on Linux
- Test with docker desktop on Mac
- Test with podman desktop on Windows
- Test with docker desktop on Windows
- Add goreleaser config to release to Github Releases
- Add Github actions workflow to build, test and release

## Questions

- Should we add colour support to output?
- Should we support Chef? (No real need as they have test-kitchen?)
- Should we support InSpec? (I think probably yes, as it's pretty awesome)
- Should we run testinfra/inspec inside the container or from outside? (Nice to not need all the python/ruby environments/packages on the host)
- Should we support Windows? (Yeah...probably)
- Should we support other shells than bash? Not sure we need to test with alpine as most people don't run their servers on alpine, this ain't Kubernetes
