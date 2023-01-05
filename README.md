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

```text
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
............

Total Duration: 0.016s
Count: 12, Failed: 0, Skipped: 0

   • verifying container rolecule-sshd-ubuntu-22.04 with goss
............

Total Duration: 0.015s
Count: 12, Failed: 0, Skipped: 0

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

## Provisioners

### Ansible

The default for the ansible provisioner is the equivalent to setting the following in the `rolecule.yml` config file:

```yaml
provisioner:
  name: ansible
  command: ansible-playbook
  playbook: playbook.yml
  args:
    - --connection
    - local
    - --inventory
    - localhost,
    - tests/playbook.yml
  extra_args: []
  env:
    ANSIBLE_ROLES_PATH: .
    ANSIBLE_NOCOWS: True
```

If you want to add extra environment variables you can just add them to the `env` map, e.g.:

```yaml
provisioner:
  name: ansible
  env:
    ANSIBLE_NOCOLOR: True
```

If you want to run a completely different ansible command, you can override the command and all the args with the `command` and `args` keys respectively, but if you just want to add other args like `--diff` or `--verbose`, add them to the `extra_args` array, e.g.:

```yaml
provisioner:
  name: ansible
  extra_args:
    - --diff
    - --verbose
```

## Verifiers

### Goss

The default goss configuration when you only specify the name, if equivalent to this:

```yaml
verifier:
  name: goss
  testfile: goss.yaml
  extra_args: []
```

This will execute:

```text
goss --gossfile tests/goss.yaml validate
```

If you want to customise how goss validate is executed, you can change the gossfile and add extra arguments to the validate subcommand, e.g. with this in your `rolecule.yml`:

```yaml
verifier:
  name: goss
  testfile: goss_tests.yaml
  extra_args:
    - --format
    - tap
```

It will execute:

```text
goss --gossfile tests/goss_tests.yaml validate --format tap
```

### FAQ

**How do I get this working on macOS?**

You'll need to make sure you create a podman machine with your home directory mounted for volume mounts to work, e.g.:

```text
» podman machine init --now --rootful -v $HOME:$HOME
```

**How do I create a suitable container image for this?**

You can use the `Containerfile`/`Dockerfile` files in the testing directory to build suitable images:

```text
» podman build -t rockylinux-systemd:9.1 -f testing/ansible/rockylinux-9.1-systemd.Containerfile .

» podman build -t ubuntu-systemd:22.04 -f testing/ansible/ubuntu-22.04-systemd.Containerfile .
```

## TODO

- ~~Test with podman on Mac~~
- ~~Test docker on Linux~~
- Make provisioner output **unbuffered**
- Support installing ansible collections
- Support testinfra verifier
- Support scenarios, making it possible to test a role with different tags
- ~~Support using custom provisioner command/args/env vars from rolecule.yml~~
- ~~Support using custom verifier command/args/env vars from rolecule.yml~~/%s
- Test converging with puppet apply
- Implement `rolecule init` to generate a rolecule.yml file (use current directory structure to determine configuration management provisioner)
- ~~Implement `rolecule list` subcommand to list all running containers~~
- ~~Write some tests :/~~
- Document what is required for a container image
- ~~Test with docker on Linux~~
- ~~Test with docker desktop on Mac~~
- ~~Test with podman desktop on Windows~~
- ~~Test with docker desktop on Windows~~
- Add goreleaser config to release to Github Releases
- Add Github actions workflow to build, test and release

## Questions

- Should we add colour support to output?
- Should we support Chef? (No real need as they have test-kitchen?)
- Should we support InSpec? (I think probably yes, as it's pretty awesome)
- Should we run testinfra/inspec inside the container or from outside? (Nice to not need all the python/ruby environments/packages on the host)
- Should we support other shells than bash? Not sure we need to test with alpine as most people don't run their servers on alpine, this ain't Kubernetes
