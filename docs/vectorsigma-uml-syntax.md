# VectorSigma UML Syntax Documentation

VectorSigma utilizes a specific UML syntax to generate finite state machines
(FSM) from UML diagrams. This documentation outlines the key components of the
UML syntax, including title, initial state, final state, actions, guards, and
transitions.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->

## Table of Contents

- [VectorSigma UML Syntax Documentation](#vectorsigma-uml-syntax-documentation)
  - [1. Title](#1-title)
  - [2. Initial State](#2-initial-state)
  - [3. Final State](#3-final-state)
  - [4. Actions](#4-actions)
    - [4.1 Good Practices for Naming](#41-good-practices-for-naming)
  - [5. Guards](#5-guards)
    - [5.1 Guarded vs. Unguarded Transitions](#51-guarded-vs-unguarded-transitions)
      - [Example of Guarded and Unguarded Transitions](#example-of-guarded-and-unguarded-transitions)
    - [5.2 Guarded Action Transitions](#52-guarded-action-transitions)
      - [Action Parameters](#action-parameters)
      - [Execution Flow](#execution-flow)
      - [Important Caveats](#important-caveats)
      - [Syntax](#syntax)
  - [6. Transitions](#6-transitions)
  - [7. Composite States](#7-composite-states)
    - [7.1 Defining Composite States](#71-defining-composite-states)
  - [8. Notes](#8-notes)

<!-- markdown-toc end -->

## 1. Title

The `title` keyword in the UML syntax is used to specify the title of the state
machine diagram. This title serves a dual purpose: it not only provides a clear
and concise label that describes the purpose or identity of the diagram but is
also utilized by VectorSigma to name the corresponding state machine struct in
the generated code.

For example:

```plaintext
title VectorSigma
```

In this case, the title "VectorSigma" indicates that the diagram pertains to the
VectorSigma application and will be reflected in the generated state machine
struct.

## 2. Initial State

The initial state is represented by the notation `[*]`. It indicates the
starting point of the state machine. For example:

```plantuml
[*] --> StateA
```

This indicates that the state machine begins in the `StateA`.

## 3. Final State

Final states are also represented by the notation `[*]`. The key difference
between the initial state and final state notation is that the initial state
will always be the first state defined in the diagram, while the final state is
always the target of a transition. For example:

```plantuml
StateB -[dotted]-> [*]: IsError
```

In this case, the transition from `StateB` to the final state occurs if the
guard condition `IsError` evaluates to true. This indicates that the state
machine can terminate upon encountering an error during the `StateB` phase.

## 4. Actions

Actions are the operations that run in a state. In the UML syntax used by
VectorSigma, actions can optionally have a `do /` prefix. This prefix indicates
that the action is performed when entering the state. For example:

```plantuml
StateA: do / PerformActionA
```

In this case, the action `PerformActionA` is executed when the state machine
enters the `StateA`.

A state can have multiple actions by adding multiple lines with the action
syntax:

```plantuml
StateB: do / PrepareData
StateB: do / ValidateData
StateB: do / TransformData
```

In this example, when the state machine enters `StateB`, it will execute
`PrepareData`, `ValidateData`, and `TransformData` in the order they are
defined.

### 4.1 Good Practices for Naming

- **State Names**: It is a good practice to name states using the `-ing` suffix
  to indicate an ongoing process. For example, `Loading`, `Processing`, or
  `Finalizing`.

- **Action Names**: Actions should typically be named using a verb that
  describes what is being done, often prefixed with `do /`. For example,
  `do / LoadData`, `do / ProcessInput`, or `do / FinalizeOutput`.

- **Guard Names**: Guards should be prefixed with `Is` or `Has` to clearly
  indicate a condition. For example, `IsError`, `HasData`, or `IsComplete`.

## 5. Guards

Guards are conditions that must be satisfied for a transition to occur. In
VectorSigma's UML syntax, guards can be wrapped in square brackets `[ ]`. For
example:

```plantuml
StateA --> StateB: IsConditionMet      // Guard without brackets
StateA --> StateB: [ IsConditionMet ]  // Guard with brackets
```

Both transitions from `StateA` to `StateB` occur only if the guard condition
`IsConditionMet` evaluates to true. The use of square brackets is optional and
does not change the meaning of the guard; it is simply a matter of visual
preference.

### 5.1 Guarded vs. Unguarded Transitions

- **Guarded Transitions**: These transitions have a guard condition that must be
  satisfied for the transition to take place. For example:

  ```plantuml
  StateA --> StateB: IsConditionMet
  ```

  This transition will only occur if the condition `IsConditionMet` is true.

- **Unguarded Transitions**: These transitions do not have any guard conditions
  and will occur unconditionally unless a guarded transition in the current
  state evaluates to true. For example:

  ```plantuml
  StateA --> StateC
  ```

  In this case, the transition from `StateA` to `StateC` is unguarded.
  VectorSigma evaluates state transitions from top to bottom, meaning that if
  there are multiple guarded transitions, the unguarded transition will always
  be the last one considered.

#### Example of Guarded and Unguarded Transitions

Consider the following transitions from `StateA`:

```plantuml
StateA --> StateB: IsInvalid
StateA --> StateC: IsNotAuthenticated
StateA --> StateD
```

In this example:

- The transition to `StateB` will occur if `IsInvalid` is true.
- The transition to `StateC` will occur if `IsNotAuthenticated` is true.
- The transition to `StateD` is unguarded and will occur only if neither of the
  previous conditions are met. If both `IsInvalid` and `IsNotAuthenticated` are
  false, the state machine will transition to `StateD`.

### 5.2 Guarded Action Transitions

VectorSigma supports a special type of transition called "guarded action
transitions" that combine both a guard condition and an action to be executed
when that guard condition is true. These transitions use the double colon (`::`)
as a separator between the guard and the action.

#### Syntax

```plantuml
StateA --> StateB: IsConditionMet::PerformAction
```

In this example, the transition from `StateA` to `StateB` will occur if the
guard condition `IsConditionMet` evaluates to true. Additionally, the action
`PerformAction` will execute as part of the transition.

#### Execution Flow

The execution flow for a guarded action transition is as follows:

1. The guard condition is evaluated.
2. If the guard condition is false, the transition is not taken, and the next
   guard condition is evaluated.
3. If the guard condition is true, the associated action is executed.
4. If the action executes successfully, the transition to the target state
   occurs.
5. If the action fails (returns an error), the state machine will immediately
   transition to the final state with an error status.

#### Important Caveats

**WARNING**: A critical caveat with guarded action transitions is that if the
action fails, the state machine will ALWAYS transition directly to the final
state. This is different from regular actions within states, which can be
handled with error transitions.

Due to this behavior, guarded action transitions should be used with caution and
are best suited for:

- Extended state status updates that cannot have system failures
- Simple operations with minimal failure risk
- Cases where immediate termination is an acceptable response to a failure

#### Action Parameters

Guarded actions can include parameters just like normal actions. These
parameters are passed to the action function when it is executed:

```plantuml
StateA --> StateB: IsConditionMet::PerformAction(param1, param2)
```

In this example, if `IsConditionMet` is true, the `PerformAction` will be called
with the parameters "param1" and "param2". This allows for more flexible and
context-specific actions to be executed as part of guarded transitions.

## 6. Transitions

Transitions define the movement from one state to another. In the UML syntax,
transitions are represented by arrows. VectorSigma supports two types of arrow
decorations:

- `-[dotted]->`: Represents a transition that may indicate a less significant or
  optional path.
- `-[bold]->`: Represents a default or significant transition.

For example:

```plantuml
StateB -[dotted]-> [*]: IsError
StateA -[bold]-> StateC
```

In this case, the transition from `StateA` to `StateC` is a default transition,
while the transition to the final state `[ * ]` from `StateB` is marked as
optional (dotted).

**Note:** The arrow decorations (`dotted` and `bold`) are ignored by VectorSigma
during the FSM generation process. They serve only as visual indicators in the
UML diagram.

## 7. Composite States

Composite states (also known as hierarchical or nested states) allow you to
organize states into hierarchies, which can simplify complex state machines by
grouping related states together. VectorSigma supports the definition and
processing of composite states in UML diagrams.

### 7.1 Defining Composite States

To define a composite state in the UML syntax, you use the `state` keyword
followed by the state name and curly braces to enclose the nested states:

```plantuml
state CompositeState {
  [*] --> NestedState1

  NestedState1: do / PerformNestedAction1
  NestedState1 --> NestedState2: IsNestedConditionMet

  NestedState2: do / PerformNestedAction2
  NestedState2 --> [*]
}
CompositeState -[dotted]-> [*]: IsError
CompositeState -[bold]-> StateA
```

In this example, `CompositeState` contains two nested states: `NestedState1` and
`NestedState2`.

## 8. Notes

The UML diagram may contain notes that provide additional context or
explanations. In the provided UML, notes are included as follows:

```plantuml
note left of StateA
  Some description of StateA
  and its actions and guards
end note
```

Currently, these notes are ignored by VectorSigma during the FSM generation.
However, there is consideration for incorporating them as function documentation
in future versions.
