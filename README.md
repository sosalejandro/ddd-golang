# DDD Golang Project

This project is an implementation of Domain-Driven Design (DDD) principles using Golang. It aims to provide a structured approach to building complex software systems by focusing on the core domain and domain logic.

## Table of Contents

- [DDD Golang Project](#ddd-golang-project)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Getting Started](#getting-started)
  - [TODO](#todo)
    - [Internal TODO](#internal-todo)
  - [Contributing](#contributing)
  - [License](#license)

## Introduction

This project demonstrates how to apply DDD principles in a Golang application. It includes examples of Aggregates, Entities, Value Objects, and Domain Events.

## Getting Started

To get started with this project, clone the repository and install the necessary dependencies.

```bash
git clone https://github.com/sosalejandro/ddd-golang.git
cd ddd-golang
go mod tidy
```


## TODO

- [x] Aggregate Root - See the implementation in the [/pkg/aggregate](./pkg/aggregate/README.md) folder
- [ ] Improvements for standardization and domain errors
- [x] Event Manager (Serializer for an Event Store)
- [ ] Open Telemetry [OTEL](https://opentelemetry.io/) first-class support. 


### Internal TODO

- [x] Implement cassandra methods for Load and SaveSnapshot
- [x] Implement InMemory Repository
- [x] Create PoC for InMemory repository implementation
- [ ] ~~Create PoC for Cassandra's repository implementation~~
- [ ] Update `AggregateRoot` documentation for its implementation
- [ ] Generate ADR for `AggregateRoot` refactor
- [x] Implement Moq library
- [x] Implement generic repository tests (mock AggregateRepositoryInterface)
- [ ] ~~Implement cassandra mock tests~~
- [x] Implement mocked tests for factory methods under the `factory` folder
- [x] Refactor Golang modules
- [ ] Generate Repository documentation
- [ ] Generate EventManager documentation

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.