droot
=====

Droot is a super-easy container chrooting a docker image without docker.

## Overview

[Docker](https://www.docker.com) has a powerful concept about an application deployment process, that is Build, Ship, Run. But there are many cases that docker runtime is beyond our current capabilities.

Droot provides a simpler container runtime without annoying Linux Namespaces. Droot depends on traditional Linux functions such as chroot(2), Linux capabilities(7) and a bind mount. `droot` helps you to chroot a container image built by `docker` and to import/export container images on Amazon S3.

![droot concept](http://cdn-ak.f.st-hatena.com/images/fotolife/y/y_uuki/20151129/20151129193210.png?1448793174)

## Requirements

- Docker (`droot push` only depends on it)
- Linux (`droot run` and `droot umount` only supports it)

## Installation

## Usage

```bash
$ droot push --to s3://example.com/dockerfiles/app.tar.gz dockerfiles/app
```

```bash
$ droot pull --dest /var/containers/app --src s3://example.com/dockerfiles/app.tar.gz
```

```bash
$ sudo droot run --bind /var/log/ --root /var/containers/app command
```

```bash
$ sudo droot umount --root /var/containers/app
```

### How to set your AWS credentials

droot push/pull subcommands support the following methods to set your AWS credentials.

1. an IAM instance profile. http://docs.aws.amazon.com/codedeploy/latest/userguide/how-to-create-iam-instance-profile.html
1. Environment variables.

```bash
$ export AWS_ACCESS_KEY_ID=********
$ export AWS_SECRET_ACCESS_KEY=********
$ export AWS_REGION=********
```

1. `~/.aws/credentials` [a standard to manage credentials in the AWS SDKs](http://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs)

## Roodmap

- `rm` command to clean container environment
- `rmi` command to clean a image on S3
- `pull` command with the rsync option
- `push/pull` other compression algorithms
- image versioning
- `pull` from docker registry
- drivers except Amazon S3
- `run` reads `.docekrenv`, `.dockerinit`
- reduce fork&exec

## Contribution

1. Fork ([https://github.com/yuuki1/droot/fork](https://github.com/yuuki1/droot/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `make test` command and confirm that it passes
1. Create a new Pull Request

## Author

[y_uuki](https://github.com/yuuki1)
