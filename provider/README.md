# Golem Provider Prototype

This project demonstrates how to use the Golem provider to rent compute resources to decentralized requestors on the Golem Network.

## Prerequisites

First, build the Golem provider plugin `ya-runtime-salad` as described in the [ya-runtime-salad/README.md](./ya-runtime-salad/README.md).

Then, ensure you have the following installed:

- [Docker](https://www.docker.com/get-started)

## Building

Build the Docker image:

```bash
docker build -t golem-provider:latest .
```

## Running

The following example will demonstrate running three different providers with different hardware configurations.

First, create separate wallets for each provider using Ethers.js:

```js
const { ethers } = require("ethers");

const wallet = ethers.Wallet.createRandom();

console.log("Address:", wallet.address);
console.log("Mnemonic:", wallet.mnemonic.phrase);
console.log("Private Key:", wallet.privateKey);
```

Save the mnemonics if you would like to reuse these example wallets in the future. To regenerate the private key from an existing mnemonic using Ether.js:

```js
const { ethers } = require("ethers");

const mnemonic = "your mnemonic phrase here";
const wallet = ethers.Wallet.fromMnemonic(mnemonic);

console.log("Address:", wallet.address);
console.log("Mnemonic:", wallet.mnemonic.phrase);
console.log("Private Key:", wallet.privateKey);
```

Next, copy the private keys into the `examples/providerX.env` files (replacing the placeholder `YAGNA_AUTOCONF_ID_SECRET` value).

Next, calculate the hourly rates for each provider based on the GPU and current Golem token value.

$$
RTX\ 3060\ Ti\ rate = (0.035 / GLM\ token\ value) / 3600
$$

$$
RTX\ 4090\ rate = (0.165 / GLM\ token\ value) / 3600
$$

Next, copy the rates into the `examples/presets-X.json` files (replacing the placeholder `golem.usage.duration_sec` value).

Finally, run the Docker containers:

```bash
docker container run --detach --cpus="2" --memory="1g" --env-file ./examples/provider1.env --rm --volume ./ya-runtime-salad/examples/template-rtx-3060-ti.json:/home/ubuntu/.config/ya-runtime-salad/template.json:ro --volume ./examples/presets-rtx-3060-ti.json:/home/ubuntu/.local/share/ya-provider/presets.json:ro golem-provider:latest run --no-interactive
```

```bash
docker container run --detach --cpus="2" --memory="1g" --env-file ./examples/provider2.env --rm --volume ./ya-runtime-salad/examples/template-rtx-4090.json:/home/ubuntu/.config/ya-runtime-salad/template.json:ro --volume ./examples/presets-rtx-4090.json:/home/ubuntu/.local/share/ya-provider/presets.json:ro golem-provider:latest run --no-interactive
```

```bash
docker container run --detach --cpus="2" --memory="1g" --env-file ./examples/provider3.env --rm --volume ./ya-runtime-salad/examples/template-rtx-4090.json:/home/ubuntu/.config/ya-runtime-salad/template.json:ro --volume ./examples/presets-rtx-4090.json:/home/ubuntu/.local/share/ya-provider/presets.json:ro golem-provider:latest run --no-interactive
```
