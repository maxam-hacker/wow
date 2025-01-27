# "Word Of Wisdom" test case

# Build and start

docker compose build

docker compose up server

docker compose up client

docker compose up tester

# Implementation details

## Transport

It is used `epoll` functionality for the server implementation.

## Prof Of Work

The project uses Hashcash algorith as PoW.
It is simple for understanding, implementation for the client and verification for the server.

## How to improve the project

Use a database as a storage of the book content.

If it is planning to use this service in a distributed architecture there will be a need using a distributed cache system (like Redis or etcd).

More metcis and control algorithms based on it.

