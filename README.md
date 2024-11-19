# DDD Golang Project

This project is an implementation of Domain-Driven Design (DDD) principles using Golang. It aims to provide a structured approach to building complex software systems by focusing on the core domain and domain logic.

## Table of Contents

- [DDD Golang Project](#ddd-golang-project)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Getting Started](#getting-started)
  - [TODO](#todo)
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
- [ ] Event Manager (Serializer for an Event Store)
- [ ] Open Telemetry [OTEL](https://opentelemetry.io/) first-class support. 

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.