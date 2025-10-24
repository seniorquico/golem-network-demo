# Golem Provider Prototype

This project demonstrates how to use the Golem provider to rent compute resources to decentralized requestors on the Golem Network.

## Prerequisites

First, build the Golem provider plugin `ya-runtime-salad` as described in the [ya-runtime-salad/README.md](./ya-runtime-salad/README.md).

Then, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)

## Building

Build the Docker image:

```bash
docker build -t golem-provider-prototype .
```
