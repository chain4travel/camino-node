<div align="center">
  <img src="resources/camino-logo.png?raw=true">
</div>

---

Node implementation for the [Camino](https://camino.foundation) network -
a blockchains platform for the touristic market.

## Installation

Camino is an quite lightweight protocol which requires only modest minimum computer requirements.

The minimum recommended hardware specification for nodes connected to Mainnet is:

- CPU: Equivalent of 8 AWS vCPU
- RAM: 16 GiB
- Storage: 512 GiB
- OS: Ubuntu 18.04/20.04 or macOS >= 10.15 (Catalina)
- Network: Reliable IPv4 or IPv6 network connection, with an open public port.

If you plan to build Camino-Node from source, you will also need the following software:

- [Go](https://golang.org/doc/install) version >= 1.17.9
- [gcc](https://gcc.gnu.org/)
- g++

### Native Install

Clone the Camino-Node repository:

```sh
git clone git@github.com:chain4travel/camino-node.git
cd camino-node
```

This will clone and checkout to `chain4travel` branch.

#### Building the Camino Node Executable

Build camino-node using the build script:

```sh
./scripts/build.sh
```

The Camino binary, named `camino-node`, is in the `build` directory.

### Binary Install

Download the [latest build](https://github.com/chain4travel/camino-node/releases/latest) for your operating system and architecture.

The Camino binary to be executed is named `camino-node`.

### Docker Install

Make sure docker is installed on the machine - so commands like `docker run` etc. are available.

Building the docker image of latest camino-node branch can be done by running:

```sh
./scripts/build_local_image.sh
```

To check the built image, run:

```sh
docker image ls
```

The image should be tagged as `chain4travel/camino-node:xxxxxxxx`, where `xxxxxxxx` is the shortened commit of the Camino source it was built from. To run the Camino node, run:

```sh
docker run -ti -p 9650:9650 -p 9651:9651 chain4travel/camino-node:xxxxxxxx /camino-node/build/camino-node
```

## Running Camino

### Connecting to Columbus Testnet

To connect to the Columbus Testnet, run:

```sh
./build/camino-node --network-id=columbus
```

You should see some pretty ASCII art and log messages.

You can use `Ctrl+C` to kill the node.

### Connecting to Camino Mainnet

Not yet launched

### Creating a Local Testnet

See [this tutorial.](https://docs.camino.foundation/developer/build/create-a-local-test-network/)

## Bootstrapping

A node needs to catch up to the latest network state before it can participate in consensus and serve API calls.

A node will not [report healthy](https://docs.camino.foundation/developer/apis/camino-node-apis/health) until it is done bootstrapping.

## Generating Code

Camino-Node uses multiple tools to generate efficient and boilerplate code.

## Supported Platforms

Camino-Node can run on different platforms, with different support tiers:

- **Tier 1**: Fully supported by the maintainers, guaranteed to pass all tests including e2e and stress tests.
- **Tier 2**: Passes all unit and integration tests but not necessarily e2e tests.
- **Tier 3**: Builds but lightly tested (or not), considered _experimental_.
- **Not supported**: May not build and not tested, considered _unsafe_. To be supported in the future.

The following table lists currently supported platforms and their corresponding
Camino-Node support tiers:

| Architecture | Operating system | Support tier  |
| :----------: | :--------------: | :-----------: |
|    amd64     |      Linux       |       1       |
|    arm64     |      Linux       |       2       |
|    amd64     |      Darwin      |       2       |
|    amd64     |     Windows      |       3       |
|     arm      |      Linux       | Not supported |
|     i386     |      Linux       | Not supported |
|    arm64     |      Darwin      | Not supported |

To officially support a new platform, one must satisfy the following requirements:

| Camino-Node continuous integration    | Tier 1  | Tier 2  | Tier 3  |
| ---------------------------------- | :-----: | :-----: | :-----: |
| Build passes                       | &check; | &check; | &check; |
| Unit and integration tests pass    | &check; | &check; |         |
| End-to-end and stress tests pass   | &check; |         |         |

## Security Bugs

**We and our community welcome responsible disclosures.**

If you've discovered a security vulnerability, please report it to us via [discord](https://discord.gg/K5THjAweFB). Valid reports will be eligible for a reward (terms and conditions apply).
