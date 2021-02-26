# Realworld Service

This is the Realworld service

Generated with

```
micro new --namespace=com.example --alias=realworld --type=service realworld-example-app
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.example.service.realworld
- Type: service
- Alias: realworld

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./realworld-service
```

Build a docker image
```
make docker
```