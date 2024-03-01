# WASAPhoto 
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/license/MIT)

Exam project for Web and Software Architecture course, year 2023/2024.

## Project structure
Structure from [Fantastic coffee (decaffeinated)](https://github.com/sapienzaapps/fantastic-coffee-decaffeinated.git):

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is an example of a daemon for checking the health of servers daemons; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the documentation (usually, for APIs, this means an OpenAPI file)
* `service/` has all packages for implementing project-specific functionalities
	* `service/api` contains an example of an API server
	* `service/globaltime` contains a wrapper package for `time.Time` (useful in unit testing)
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is an example of a web frontend in Vue.js; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

## Deployment
### Build the images
Backend
```
docker build -t wasaphoto-backend:latest -f Dockerfile.backend .
```

Frontend
```
docker build -t wasaphoto-frontend:latest -f Dockerfile.frontend .
```

### Run the containers
Backend
```
docker run -it --rm -p 3000:3000 wasaphoto-backend:latest
```

Frontend
```
docker run -it --rm -p 8080:80 wasaphoto-frontend:latest
```
