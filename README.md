# VectorSigma

Welcome to **VectorSigma**, a powerful FSM (Finite State Machine) generator,
with a name inspired by the Transformers universe! Revered for its power to
grant knowledge and wisdom to the Cybertronians, Vector Sigma serves as a
central hub of intelligence, enabling the creation of sentient beings and
guiding them on their journeys. Similarly, our application seeks to empower
developers by streamlining the creation and management of state machines,
providing the tools and insights needed to elevate their projects to new
heights.

## Overview

VectorSigma takes a PlantUML file as input and generates a runnable state
machine in Go. With the `--init` flag, you can create a new Go module that
includes a self-contained FSM application based on your UML diagram. The
generated state machine can be easily regenerated to reflect changes in the
UML, while preserving any actions and guards you've already defined. The only
time these will be removed is if they are deleted from the UML itself.

## Features

- **FSM Generation**: Convert your PlantUML diagrams into runnable state
  machines in Go.
- **Incremental Updates**: Regenerate your state machine without losing your
  custom actions and guards.
- **Easy Initialization**: Use the `--init` flag to create a new Go module with
  a complete FSM application.

## Writing UML

For creating your UML diagrams, we recommend using
[PlantText](https://planttext.com), a user-friendly online tool that allows you
to write and visualize PlantUML code effortlessly. If you're looking to deepen
your understanding of state machines and state charts, we recommend visiting
[Statecharts.dev](https://statecharts.dev) for comprehensive resources.

## Getting Started

To get started with VectorSigma, follow these steps:

1. Install the application.
2. Create your PlantUML file defining the state machine.
3. Run the VectorSigma command with the appropriate flags to generate your FSM.

## Contributing

We welcome contributions to VectorSigma! If you have suggestions, bug reports,
or feature requests, please feel free to open an issue or submit a pull
request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file
for more details.

---

Join us on this journey to create robust state machines with VectorSigma, where
the power of Transformers meets the elegance of Go programming!
