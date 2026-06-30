Vaultsh

A read-only virtual shell engine.

Everything in this environment is virtual.
Commands operate on mounted content vaults rather than a real operating system.

Mounted vaults:

cv/          professional profile
projects/    personal projects
docs/        project documentation

Getting started:

help
tree
cd cv
cat README.txt

Useful commands:

help
pwd
ls
cd
cat
tree
grep
history
clear

Note:

This is not a Linux shell.
It is a backend engine built in Go with a virtual filesystem and command execution API.