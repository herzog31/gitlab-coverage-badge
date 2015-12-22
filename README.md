[![Build Status](https://travis-ci.org/herzog31/gitlab-coverage-badge.svg?branch=master)](https://travis-ci.org/herzog31/gitlab-coverage-badge)
[![Docker Hub](https://img.shields.io/docker/pulls/herzog31/gitlab-coverage-badge.svg)](https://hub.docker.com/r/herzog31/gitlab-coverage-badge)
[![Release](https://img.shields.io/github/release/herzog31/gitlab-coverage-badge.svg)](https://github.com/herzog31/gitlab-coverage-badge/releases)
[![Go](https://img.shields.io/badge/Go-1.5.1-blue.svg)](https://golang.org/)

# Coverage Badge for Gitlab
**Currently not working with Gitlab CE 8.3!**

# Installation
Run the container using the following command and adapt the environment variables and port as necessary.

```
docker run -d -e TOKEN=YOURTOKEN -e GITLAB_HOST=https://gitlab.example.org -p 8080:8080 herzog31/gitlab-coverage-badge
```

# Usage
If your project in Gitlab is named `name/project`, then you find its coverage badge at:

```
http://your.docker.host/name/project.svg
```

Currently, only output as SVG is supported.