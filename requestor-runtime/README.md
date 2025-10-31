# Golem Requestor Prototype

This project demonstrates how to use the Golem requestor to rent compute resources from decentralized providers on the Golem Network.

## Prerequisites

Ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)

## Building

Build the Docker image:

```bash
docker build -t golem-requestor:latest .
```

## Running

The following example will demonstrate running a requestor.

First, create a wallet for the requestor using Ethers.js:

```js
const { ethers } = require("ethers");

const wallet = ethers.Wallet.createRandom();

console.log("Address:", wallet.address);
console.log("Mnemonic:", wallet.mnemonic.phrase);
console.log("Private Key:", wallet.privateKey);
```

Save the mnemonic if you would like to reuse this example wallet in the future. To regenerate the private key from an existing mnemonic using Ether.js:

```js
const { ethers } = require("ethers");

const mnemonic = "your mnemonic phrase here";
const wallet = ethers.Wallet.fromMnemonic(mnemonic);

console.log("Address:", wallet.address);
console.log("Mnemonic:", wallet.mnemonic.phrase);
console.log("Private Key:", wallet.privateKey);
```

Next, copy the private key into the `examples/organization1.env` file (replacing the placeholder `YAGNA_AUTOCONF_ID_SECRET` value).

Finally, run the Docker container:

```bash
docker container run --cpus="2" --detach --env-file ./examples/organization1.env --memory="1g" --name golem_requestor_1 --publish 7465:7465/tcp --rm golem-requestor:latest service run
```

## Funding

If you are using the testnet, you may request test tokens by running the following command:

```bash
docker container exec -it golem_requestor_1 /home/ubuntu/.local/bin/yagna payment fund --network hoodi
```

## Creating Demand

The following Golem JS SDK example demonstrates creating demand for a Salad provider:

```js
const { GolemNetwork } = await import('@golem-sdk/golem-js')

let glm = new GolemNetwork({
  api: {
    key: "6ef69d2c-e72c-4e81-a4c4-95f9120dd958",
    url: "http://127.0.0.1:7465"
  },
  payment: {
    network: "hoodi"
  }
})
await glm.connect()

let shutdown = new AbortController()
let rentalPromise = glm.oneOf({
  order: {
    demand: {
      workload: {
        runtime: {
          name: "salad"
        },
        imageTag: "golem/alpine:latest"
      }
    },
    market: {
      rentHours: 1,
      pricing: {
        model: 'linear',
        maxStartPrice: 0.01,
        maxCpuPerHourPrice: 0.01,
        maxEnvPerHourPrice: 2.8312571
      },
    }
  },
  signalOrTimeout: shutdown.signal
})
```
