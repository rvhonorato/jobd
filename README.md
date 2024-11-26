# `jobd` (Job Daemon)

![GitHub License](https://img.shields.io/github/license/rvhonorato/jobd)
[![ci](https://github.com/rvhonorato/jobd/actions/workflows/ci.yml/badge.svg)](https://github.com/rvhonorato/jobd/actions/workflows/ci.yml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/ae50eb8e1303415f981ec755f0b8a28f)](https://app.codacy.com/gh/rvhonorato/jobd/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/ae50eb8e1303415f981ec755f0b8a28f)](https://app.codacy.com/gh/rvhonorato/jobd/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)

`jobd` is a lightweight Golang application designed to facilitate interaction with
research software through REST APIs. It is specifically engineered to be deployed
in multi-stage Docker builds, providing a flexible and portable solution for job
management and file transfer.

## Key Features

Implements two primary REST API endpoints:

- `/upload`: Allows backend systems or scripts to upload files to the container
- `/download`: Enables retrieval of files from the container

## Usage

The application is optimized for containerized environments,
supporting multi-stage build patterns or simple binary execution.

In both ways the `jobd` version can be passed as a build argument:

```bash
docker build --build-arg JOBD_VERSION=v0.1.0 .
```

### Multi-stage build

First stage pulls the jobd executable from a specific image
Second stage incorporates the executable into a base research container
Enables seamless integration of job management capabilities into existing research
software containers

```dockerfile
ARG JOBD_VERSION=latest
FROM ghcr.io/rvhonorato/jobd:${JOBD_VERSION} AS jobd

FROM ghcr.io/haddocking/arctic3d:v0.5.1 AS base

WORKDIR /data
COPY --from=jobd /path/to/jobd /bin/jobd

ENTRYPOINT [ "/bin/jobd" ]
```

### Binary execution

```dockerfile
FROM ghcr.io/haddocking/arctic3d:v0.5.1
WORKDIR /data

ARG JOBD_VERSION=v0.1.0
ARG JOBD_ARCH=linux_386

# Download and extract jobd binary
ADD https://github.com/rvhonorato/jobd/releases/download/${JOBD_VERSION}/jobd_${JOBD_VERSION}_${JOBD_ARCH}.tar.gz /tmp/
RUN tar -xzf /tmp/jobd_${JOBD_VERSION}_${JOBD_ARCH}.tar.gz -C /bin/ \
    && chmod +x /bin/jobd \
    && rm /tmp/jobd_${JOBD_VERSION}_${JOBD_ARCH}.tar.gz

ENTRYPOINT [ "/bin/jobd" ]
```

## Use Cases

- Microservice-based job submission and file handling
- Simplified API interfaces for research software workflows

## API description

Soon!

## Technical Characteristics

- Written in Golang for performance and simplicity
- REST API-based communication
- Lightweight and containerization-friendly
- Designed for research and scientific computing environments

The application provides a standardized, portable mechanism for programmatic file and job interactions across different research software platforms.
