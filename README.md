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

### `/api/upload`

The `/api/upload` endpoint allows users to submit a job for processing in the queue system.

## Request Details

### HTTP Method

`POST`

### Request Body

The request body should be a JSON object representing a Job, with the following key properties:

| Field    | Type    | Required | Description                           |
| -------- | ------- | -------- | ------------------------------------- |
| `ID`     | string  | Optional | Unique identifier for the job         |
| `Input`  | string  | Required | Base64 encoded .zip file containing:  |
| `Slurml` | boolean | Optional | Flag to indicate Slurm job submission |

### Input Zip File Structure

The uploaded zip file must contain:

- `run.sh`: Executable shell script that defines the application command
- Input files required by the application

Soon!

### `/api/get:id`

The `/api/get/:id` endpoint allows users to retrieve the status and details of a previously submitted job.

### Request Details

### HTTP Method

`GET`

### URL Parameters

| Parameter | Type   | Required | Description                              |
| --------- | ------ | -------- | ---------------------------------------- |
| `id`      | string | Required | Unique identifier of the job to retrieve |

## Response

### Successful Responses

#### Job Not Yet Completed

- **Status Code**: `206 Partial Content`
  - Returned when job is in a partial/incomplete state
- **Body**: Job object with limited details

#### Job Completed

- **Status Code**: `200 OK`
  - Returned when job has finished processing
- **Body**: Job object with final status

### Response Object

The returned job object will have the following key characteristics:

- Job ID preserved
- Status of the job
- Final result/output - also a base64 encoded zip file
- **Note**: `Input` and `Path` fields are deliberately cleared before response

### Possible Job Statuses

- Pending
- Processing
- Partial
- Completed
- Failed

### Error Responses

- **404 Not Found**: Job ID does not exist
- **500 Internal Server Error**: Server-side processing error

## Technical Characteristics

- Written in Golang for performance and simplicity
- REST API-based communication
- Lightweight and containerization-friendly
- Designed for research and scientific computing environments

The application provides a standardized, portable mechanism for programmatic file and job interactions across different research software platforms.
