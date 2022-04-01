# WoWPoW

Custom tcp-server with pow challenge–response DDOS protection.
Server has semi-custom interaction protocol, based on [protocol buffers](https://developers.google.com/protocol-buffers).
Porotobuf used because it is quick binary protocol.

## Proof of work challenge–response algorithm

Was chosed [hashcash](https://en.wikipedia.org/wiki/Hashcash) algorithm with sha256 hashing function.
Reasons why was chosen algorithm is simplicity of implementation any clients. Clients with arm cpu arch could waste time more efficiently and acceptance could be harder.

Community developed many opensource libraries to produce work using hashcash algorithm. 

In additionally hashchash is most known PoW algorithm ad it uses in many different cryptocurrency miners (in particular [Bitcoin](https://en.bitcoin.it/wiki/Hashcash) and many others).

Sha256 chosed because it is the most complexity hash function. Scrypt is KDF, originally and semantically this function designed to store passwords in database.
It could be reconfigured to use sha1 hash function.

Was solved two corner problems of algorithm:
- Pre-generation of challenges fixed by very short expiration time (2 minutes and could be configured by env param `HASHCASH_CHALLENGE_EXP_DURATION=5m`)
- We are 100% sure, that all challenges produced by our server. It has extension field with generate secure sha256 hash. It contain resource, rand, timestamp and secret key which known only on server. Because of this we don't need to store any data in any database. Example of generated hashcash: `1:4:1648837394:127.0.0.1:49948:6e6c7d2af33b7db1cf92b9875852cc127bc91e64f06dada45851da4c123902e5:MzUwMjkyOTY4ODgzODcwOTgxOA==:MA==`

## Simple run

Run command to start docker-compose

```shell
make start
```

## Tests

```shell
make test
```

or for more detailed coverage

```shell
make coverage
```

## Dependencies

- [Go](https://go.dev) server and client implementation
- [buf](https://buf.build/) alternate to the `protoc` protobuf generator
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Install

```shell
make install
```

### Generate

#### Generate mocks

In case if you changed interfaces, you should generate new mocks and fix tests.

```shell
make gen
```

#### Generate proto files

This action is required only in case when you edit `api/proto/*.proto` files.

```shell
make genProto
```