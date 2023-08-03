# ZKP Auth protocol

A proof of concept application implementing the Chaum–Pedersen Protocol as an authentication service.

See [Cryptography: An Introduction (3rd Edition) Nigel Smart](https://www.cs.umd.edu/~waa/414-F11/IntroToCrypto.pdf) for more details about the protocol.

## Implementation

The protocol is implemented as simple server and client applications. The server functions as the "Verifier", whereas the client functions as the "Prover".
The main implementation code is in the [zkp_auth](./zkp_auth) directory. Simple wrappers are also provided in the [server](./server) and [client](./client) directories.

The protobuf is included in [auth.proto](./auth.proto).

### Server

In order to start the server simply run 

```
go run ./server/server.go
```

This stands up a server listening on port `50051` providing three RPC methods defined in the protobuf, namely - `Register`, `CreateAuthenticationChallenge` and `VerifyAuthentication`.

### Client

The [client](./zkp_auth/client.go) provides two simple functions `Register` and `Login` making it easy to integrate into any client-side code. There is an example integration in the [wrapper code](./client/client.go).

### Docker

Both the server and the client codes are dockerised using the corresponding docker files: [Dockerfile-server](./Dockerfile-server) and [Dockerfile-client](./Dockerfile-client).

A [docker-compose](./docker-compose.yml) is also provided for convenience. Running ```docker-compose up --build``` stands up the server and the client in separate docker containers, and runs the test client code. 

The test client code includes two tests: (1) Registration + successful login attempt (2) Registration + unsuccessful login attempt.

