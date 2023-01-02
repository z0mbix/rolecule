# Rolecule

## Description

`rolecule` is a simple tool to help you test your configuration management code
works as you expect, by creating systemd enabled containers with either docker or podman, then converging them with your configured provisioner (ansible by default).

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
» rolecule test
   • creating container: rolecule-rockylinux-systemd-9.1-amd64
   • 7c694fdfb0173ed2b9b86df571de11db0afd35cd011c0523372461e5eb480965
   • converging container: rolecule-rockylinux-systemd-9.1-amd64
Using /etc/ansible/ansible.cfg as config file

PLAY [test] ********************************************************************

TASK [Gathering Facts] *********************************************************
ok: [localhost]

TASK [. : package] *************************************************************
changed: [localhost] => {"changed": true, "msg": "", "rc": 0, "results": ["Installed: perl-Time-Local-2:1.300-7.el9.noarch", "Installed: perl-mro-1.23-479.el9.x86_64", "Installed: perl-libs-4:5.32.1-479.el9.x86_64", "Installed: perl-lib-0.65-479.el9.x86_64", "Installed: perl-interpreter-4:5.32.1-479.el9.x86_64", "Installed: perl-PathTools-3.78-461.el9.x86_64", "Installed: perl-File-Path-2.18-4.el9.noarch", "Installed: perl-POSIX-1.94-479.el9.x86_64", "Installed: perl-Data-Dumper-2.174-462.el9.x86_64", "Installed: perl-Mozilla-CA-20200520-6.el9.noarch", "Installed: perl-Pod-Perldoc-3.28.01-461.el9.noarch", "Installed: perl-NDBM_File-1.15-479.el9.x86_64", "Installed: perl-IO-1.43-479.el9.x86_64", "Installed: perl-Fcntl-1.13-479.el9.x86_64", "Installed: perl-Errno-1.30-479.el9.x86_64", "Installed: perl-libnet-3.13-4.el9.noarch", "Installed: perl-Text-Tabs+Wrap-2013.0523-460.el9.noarch", "Installed: perl-DynaLoader-1.47-479.el9.x86_64", "Installed: perl-B-1.80-479.el9.x86_64", "Installed: perl-Digest-MD5-2.58-4.el9.x86_64", "Installed: perl-vars-1.05-479.el9.noarch", "Installed: perl-parent-1:0.238-460.el9.noarch", "Installed: perl-subs-1.03-479.el9.noarch", "Installed: perl-Text-ParseWords-3.30-460.el9.noarch", "Installed: perl-overloading-0.02-479.el9.noarch", "Installed: perl-overload-1.31-479.el9.noarch", "Installed: perl-Exporter-5.74-461.el9.noarch", "Installed: perl-Carp-1.50-460.el9.noarch", "Installed: perl-if-0.60.800-479.el9.noarch", "Installed: perl-TermReadKey-2.38-11.el9.x86_64", "Installed: perl-Git-2.31.1-2.el9.2.noarch", "Installed: ncurses-6.2-8.20210508.el9.x86_64", "Installed: perl-Getopt-Long-1:2.52-4.el9.noarch", "Installed: perl-Scalar-List-Utils-4:1.56-461.el9.x86_64", "Installed: perl-base-2.27-479.el9.noarch", "Installed: git-core-doc-2.31.1-2.el9.2.noarch", "Installed: perl-podlators-1:4.14-460.el9.noarch", "Installed: perl-Term-Cap-1.17-460.el9.noarch", "Installed: perl-MIME-Base64-3.16-4.el9.x86_64", "Installed: perl-Term-ANSIColor-5.01-461.el9.noarch", "Installed: perl-Symbol-1.08-479.el9.noarch", "Installed: perl-Pod-Escapes-1:1.07-460.el9.noarch", "Installed: perl-Pod-Simple-1:3.42-4.el9.noarch", "Installed: perl-HTTP-Tiny-0.076-460.el9.noarch", "Installed: perl-SelectSaver-1.02-479.el9.noarch", "Installed: perl-Net-SSLeay-1.92-2.el9.x86_64", "Installed: perl-Digest-1.19-4.el9.noarch", "Installed: perl-IO-Socket-IP-0.41-5.el9.noarch", "Installed: perl-IO-Socket-SSL-2.073-1.el9.noarch", "Installed: perl-IPC-Open3-1.21-479.el9.noarch", "Installed: perl-Error-1:0.17029-7.el9.noarch", "Installed: perl-Getopt-Std-1.12-479.el9.noarch", "Installed: perl-FileHandle-2.03-479.el9.noarch", "Installed: perl-File-stat-1.09-479.el9.noarch", "Installed: perl-File-Find-1.37-479.el9.noarch", "Installed: perl-File-Basename-2.85-479.el9.noarch", "Installed: perl-URI-5.09-3.el9.noarch", "Installed: perl-constant-1.33-461.el9.noarch", "Installed: perl-Class-Struct-0.66-479.el9.noarch", "Installed: git-2.31.1-2.el9.2.x86_64", "Installed: groff-base-1.22.4-10.el9.x86_64", "Installed: perl-Encode-4:3.08-462.el9.x86_64", "Installed: perl-AutoLoader-5.74-479.el9.noarch", "Installed: perl-File-Temp-1:0.231.100-4.el9.noarch", "Installed: perl-Pod-Usage-4:2.01-4.el9.noarch", "Installed: emacs-filesystem-1:27.2-6.el9.noarch", "Installed: perl-Storable-1:3.21-460.el9.x86_64", "Installed: perl-Socket-4:2.031-4.el9.x86_64"]}

TASK [. : template] ************************************************************
--- before
+++ after: /root/.ansible/tmp/ansible-local-221haxj_dm/tmpql4yub50/test1.txt.j2
@@ -0,0 +1 @@
+test1

changed: [localhost] => {"changed": true, "checksum": "dba7673010f19a94af4345453005933fd511bea9", "dest": "/tmp/test1.txt", "gid": 0, "group": "root", "md5sum": "3e7705498e8be60520841409ebc69bc1", "mode": "0644", "owner": "root", "size": 6, "src": "/root/.ansible/tmp/ansible-tmp-1672621295.4137197-165-251753374476187/source", "state": "file", "uid": 0}

TASK [. : file] ****************************************************************
--- before
+++ after: /src/files/test2.txt
@@ -0,0 +1 @@
+test2

changed: [localhost] => {"changed": true, "checksum": "9054fbe0b622c638224d50d20824d2ff6782e308", "dest": "/tmp/test2.txt", "gid": 0, "group": "root", "md5sum": "126a8a51b9d1bbd07fddc65819a542c3", "mode": "0640", "owner": "root", "size": 6, "src": "/root/.ansible/tmp/ansible-tmp-1672621296.4018161-191-222043994690168/source", "state": "file", "uid": 0}

TASK [. : directory] ***********************************************************
--- before
+++ after
@@ -1,5 +1,5 @@
 {
-    "mode": "0755",
+    "mode": "0750",
     "path": "/tmp/simple-directory",
-    "state": "absent"
+    "state": "directory"
 }

changed: [localhost] => {"changed": true, "gid": 0, "group": "root", "mode": "0750", "owner": "root", "path": "/tmp/simple-directory", "size": 40, "state": "directory", "uid": 0}

PLAY RECAP *********************************************************************
localhost                  : ok=5    changed=4    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

   • complete
```

## Help

```shell
» ./rolecule --help
rolecule uses docker or podman to test your
configuration management roles/recipes/modules in a systemd enabled container,
then tests them with a verifier (goss/testinfra).

Usage:
  rolecule [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  converge    Run your configuration management tool to converge the configuration
  create      Create a new container(s) to test the role in
  destroy     Destroy everything
  help        Help about any command
  init        initialise the project with a nice new rolecule.yml file
  list        list the containers
  shell       get a shell in a container
  test        Create the container(s), converge them, test them, then clean up
  verify      verify your container

Flags:
  -d, --debug   enable debug output
  -h, --help    help for rolecule

Use "rolecule [command] --help" for more information about a command.
```
