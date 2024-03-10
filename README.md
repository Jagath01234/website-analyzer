# Website-Analyzer

Website analyzer using Go 1.21 

## Description

Website analyzer using Go 1.21 to analyze few basic information of a given site lite the HTML version,Title,
Is the web page given by the target URL a login page, how many internal and external links does it have etc.
However this Project currently support only for parsing and analyzing HTML only. Links calling through javascript onclick event etc are not yet supported by this project.
This project can be used as a base project to implement extensive functionality.
Prometheus Metrics, Swagger API docs and pprof are implemented in this project. 
Most parts of the business layer, gateways, utils and helpers are covered with unit tests.

A frontend client supporting the APIs of this service can be found in the  https://github.com/Jagath01234/website-analyzer-react-client repository.
## Table of Contents

- [Requirements](#Requirements)
- [Configuration](#Configuration)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)


## Requirements
- Go 1.21
- Configured GOROOT and GOPATH.
- Configured github token or made the GOPROXY to `direct` to pull the dependencies from github. (This is required since there are go dependencies which has version higher than `v1`.)

## Configuration
- Edit the `config.json` file in the project root to change the configurations as required. 
Configurations are self-explanatory.
```json
{
  "app": {
    "port": 8080
  },
  "cache": {
    "max_size": 1000,
    "prune_size": 50,
    "expiry_time_secs": 600
  },
  "api_docs": {
    "is_enabled": true
  },
  "pprof": {
    "is_enabled": true,
    "port": 6060
  },
  "metrics": {
    "port": 7070
  },
  "worker": {
    "buffer_size": 50,
    "pool_size": 10
  }
}

```

## Installation
- Navigate to project root directory using terminal.
- Download the dependencies with `go mod tidy`.
- Run with `go run main go`.

## Usage
- There are two enpoints in this service and they are 


  

