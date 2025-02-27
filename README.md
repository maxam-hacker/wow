## "Word Of Wisdom" test case

### Build and start

`docker compose build`

`docker compose up server`

`docker compose up client`

`docker compose up tester`

## Implementation details

### Transport

It is used `epoll` functionality for the server implementation and memory optimization approach to handle messages.

### Proof Of Work

The project uses Hashcash algorith as PoW.
It is simple for understanding, implementation for the client and verification for the server and hard to brutforce.
But the algo depends on calculation abilities of hardware.

### How to improve the project

Use a database as a storage of the book content.

If it is planning to use this service in a distributed architecture there will be a need using a distributed cache system (like Redis or etcd).

More metcis and control algorithms based on it.

More unit and integration tests.

`json.Marshal` is not optimize in relation to memory consumption in that case... need to think it through.

### Results

Client:
![alt text](docs/images/client.png?raw=true)

Server:
![alt text](docs/images/server.png?raw=true)

Net settings:
![alt text](docs/images/net.png?raw=true)