# `jobd` (Job Daemon)

_documentation wip_

`jobd` is a lightweight Golang application designed to facilitate interaction with research software through REST APIs. It is specifically engineered to be deployed in multi-stage Docker builds, providing a flexible and portable solution for job management and file transfer.

## Key Features
Implements two primary REST API endpoints:

- `/upload`: Allows backend systems or scripts to upload files to the container
- `/download`: Enables retrieval of files from the container


## Docker Integration
The application is optimized for containerized environments, supporting multi-stage build patterns. In the example dockerfile:

First stage pulls the jobd executable from a specific image
Second stage incorporates the executable into a base research container
Enables seamless integration of job management capabilities into existing research software containers

```dockerfile
# Stage 1: Copy jobd executable
FROM ghcr.io/rvhonorato/jobd:latest AS jobd

# Stage 2: Base research image
FROM ghcr.io/haddocking/arctic3d:v0.5.1 AS base

WORKDIR /data
COPY --from=jobd /path/to/jobd /bin/jobd

ENTRYPOINT [ "/bin/jobd" ]
```

## Use Cases

- Microservice-based job submission and file handling
- Simplified API interfaces for research software workflows

## Technical Characteristics

- Written in Golang for performance and simplicity
- REST API-based communication
- Lightweight and containerization-friendly
- Designed for research and scientific computing environments

The application provides a standardized, portable mechanism for programmatic file and job interactions across different research software platforms.
