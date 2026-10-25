Implementation Plan — Data-Driven Combat System in Go

1. Objective

Implement a turn-based combat system for an existing solo RPG CLI written in Go.

The system must be configurable through YAML files, following the project’s existing data-driven approach used for cards and oracles.

The goal is to make combat rules, character classes, available actions, action costs, effects, and responses configurable without hardcoding class-specific gameplay into the engine.

Important: This is an existing project. Analyze the current implementation before making changes. Reuse existing structures, conventions, persistence mechanisms, YAML loading, and CLI patterns wherever appropriate. Do not redesign unrelated parts of the application.

This is an implementation task, not a request to produce another proposal or ask the user to make additional design decisions.

⸻

2. Architectural Principles

Separate the system into three conceptual layers:

1. Combat rules: Define the general combat flow and resolution behavior.
2. Action definitions: Define reusable actions, their costs, effects, and responses.
3. Class definitions: Define action availability, action points, and class-specific gameplay.

The Go engine is responsible for loading and validating definitions, maintaining combat state, executing actions, resolving results, and applying effects.

YAML defines the rules within the capabilities supported by the engine. It must not become an unrestricted programming language.

The engine must not contain hardcoded behavior specific to classes such as Assassin, Warrior, Mage, or Boss.

⸻

3. Repository Analysis Before Implementation

Before changing code:

* Inspect the existing project structure and coding conventions.
* Identify how YAML files are located, loaded, decoded, and validated.
* Inspect existing card models and their persistence/export mechanisms.
* Identify how game state, characters, items, and statistics are represented.
* Inspect the CLI interaction patterns and manual input mechanisms.
* Identify existing dice/oracle resolution functionality that can be reused.
* Inspect the current test framework and test conventions.

Use these findings to adapt the proposed design to the project.

Do not introduce duplicate abstractions when an existing project component already provides the required functionality.

Do not modify unrelated systems.

⸻

4. Proposed Data Organization

Adapt the directory names to existing project conventions.

data/
├── rules/
│   └── combat.yaml
├── actions/
│   ├── attack.yaml
│   ├── defend.yaml
│   └── counter.yaml
└── classes/
    ├── assassin.yaml
    ├── warrior.yaml
    ├── mage.yaml
    └── goblin.yaml

Responsibilities

File category	Responsibility
Combat rules	General turn flow and action-resolution behavior
Action definitions	Action cost, effects, uses, critical outcomes, and responses
Class definitions	Action points and actions available to a combatant

These are logical responsibilities. The implementation may reuse existing card storage or YAML organization rather than creating this exact directory structure.

Action IDs must be stable identifiers that can be referenced from other YAML files.

⸻

5. Combat Rules YAML

Create a general combat configuration that describes how the combat engine behaves.

Proposed initial structure:

turn:
  action_points: class
  allow_pass: true
resolution:
  responses: automatic
  max_response_depth: 10

Expected semantics:

* action_points: class: The current combatant’s action-point budget comes from its class definition.
* allow_pass: true: A combatant can end its turn before spending all available action points.
* responses: automatic: Responses configured by actions are executed automatically.
* max_response_depth: Limits response-chain depth to prevent runaway execution.

Treat this structure as an initial schema, not a reason to introduce unnecessary configuration infrastructure.

Validate all supported values during loading.

If the existing project has a better mechanism for configuration defaults, use it consistently.

⸻

6. Action Definitions

Actions are reusable data definitions.

They must not be implemented as separate hardcoded Go types for every individual action.

Proposed attack definition:

id: attack
cost: 1
effects:
  hit:
    damage: 3
  critical:
    damage: 6
    bonus_uses: 1
    bonus_cost: 0
response:
  action: defend
  actor: target

This is a proposed YAML schema that the implementation must support. It does not imply that these fields already exist in the codebase.

Required concepts

Cost

The normal action-point cost required to execute an action.

Uses

The number of available executions of an action.

Effects

Changes applied according to the resolution outcome.

Critical effects

Additional or alternative effects triggered by a critical result.

Response

A reference to an action executed by another combatant when the relevant response condition is satisfied.

Usage accounting

Keep action points and action uses conceptually separate.

* Action points represent a shared budget for the current turn.
* Uses represent the number of times an action is available for execution.
* An action may require both an available use and sufficient action points.
* An effect may grant an additional use.
* An effect may grant a free use that does not consume action points.

Do not implement free uses by globally changing an action’s normal cost. A free use must be distinguishable from a normal use.

Use explicit internal state to represent the resulting availability and cost rules.

⸻

7. Action Responses

Responses are references to existing actions, not a separate category of hardcoded combat behavior.

Proposed defense definition:

id: defend
cost: 0
effects:
  success:
    damage_reduction: 3
  critical:
    damage_reduction: 6

The attack definition may reference defend as its response.

Conceptually:

Player executes attack
        |
        v
Resolve attack outcome
        |
        v
Apply attack effects
        |
        v
Resolve configured response
        |
        v
Find defend action on target
        |
        v
Execute defend
        |
        v
Apply defense effects

The engine must resolve references dynamically against the appropriate combatant’s available actions.

It must not call a hardcoded Defend() function merely because the action ID is attack.

Response validation

At load time, validate that referenced action definitions exist.

At execution time, verify that the responding combatant can perform the referenced action under the applicable rules.

If an action is unavailable, apply a defined fallback behavior. The engine must not panic or silently corrupt the combat state.

Response chains must respect the configured depth limit.

⸻

8. Classes and Available Actions

A class defines its action-point budget and the actions available to its members.

Example:

id: assassin
action_points: 3
actions:
  - attack
  - defend
  - counter

Another class may use a different budget and different actions:

id: mage
action_points: 3
actions:
  - attack
  - fireball
  - shield
  - teleport

A class may represent a traditional character archetype or a gameplay-only role, such as Boss, Elite, or Monster.

The engine must not require a dedicated Go implementation for each class.

Classes should reference reusable action definitions rather than duplicating their contents.

If existing character cards already contain class-related information, integrate the class reference with that model rather than introducing a competing character representation.

⸻

9. Combat State

Separate immutable definitions from mutable combat state.

Definitions

Loaded YAML definitions describe:

* combat rules;
* class properties;
* available actions;
* costs;
* effects;
* response references.

These definitions should not be mutated during combat.

Runtime state

Combat state contains:

* participating combatants;
* current turn and active combatant;
* current action-point balance;
* action-use availability;
* current statistics and status;
* temporary effects;
* pending resolution information, when required.

Action-use state must be scoped correctly. One combatant’s temporary effects must not mutate another combatant’s definitions or availability.

Reuse existing game-state structures where appropriate.

Avoid storing transient combat state in persistent card definitions.

⸻

10. Go Engine Responsibilities

Implement the smallest set of components needed to execute the defined rules.

The logical responsibilities are:

Definition loading

Load combat rules, actions, and classes using the project’s existing YAML infrastructure.

Definition validation

Validate IDs, field values, references, and supported effect types before combat starts.

Combat initialization

Create runtime state for the participants and initialize their action-point balances and available actions.

Action availability

Determine whether an action can be executed by checking:

* action availability for the combatant;
* remaining uses;
* sufficient action points;
* applicable state restrictions.

Action execution

The execution process should:

1. Validate the action.
2. Reserve or consume the applicable cost and use.
3. Obtain the resolution outcome.
4. Apply the corresponding effects.
5. Resolve configured responses.
6. Verify combat termination conditions.
7. Return the updated state and relevant events to the CLI.

The implementation must handle failures without leaving partially applied state where practical.

Manual resolution

Reuse existing dice/oracle infrastructure when suitable.

If the player must provide a roll result manually, the engine should expose the information needed for that interaction and receive the result through the existing CLI conventions.

Do not introduce an LLM dependency.

⸻

11. Critical Hits and Free Attacks

Implement the critical effect mechanism using generic effect handling.

Example:

effects:
  critical:
    damage: 6
    bonus_uses: 1
    bonus_cost: 0

Expected behavior:

1. The attack resolves as a critical.
2. The configured critical damage is applied.
3. One additional attack use is granted.
4. That additional use is free.
5. The normal cost of future attacks remains unchanged.
6. Each new attack follows its configured response rules.

The implementation must distinguish a free use from an ordinary use. It must not permanently reduce the attack cost to zero.

Add tests covering multiple critical hits and the resulting action-use balance.

⸻

12. CLI Integration

Integrate combat into the existing game flow.

The CLI should:

* display the active combatant;
* display current statistics and action points;
* display available actions;
* reject unavailable or unaffordable actions;
* request the required manual resolution input;
* display effects and responses;
* allow the player to pass the turn;
* recognize victory, defeat, and other configured combat termination conditions.

Do not create a separate CLI framework if the project already has one.

Keep presentation and input handling separate from rule execution wherever consistent with the existing architecture.

⸻

13. Testing Requirements

Write unit tests for the rule engine and integration tests for the complete combat flow.

At minimum, cover the following cases.

Action points

* An affordable action executes successfully.
* An unaffordable action is rejected.
* The action-point balance is updated correctly.
* Passing a turn works even when points remain.

Action uses

* An action cannot be used when no uses remain.
* A critical grants the configured additional use.
* A free use does not consume action points.
* A normal use still consumes its regular cost.

Effects

* Normal hits apply the correct effects.
* Critical hits apply their configured effects.
* Defense modifies damage according to the configured rules.

Responses

* An attack triggers its configured defense.
* A critical defense can trigger a configured counterattack.
* An unavailable response is handled safely.
* Response chains terminate correctly.
* Cyclic response configurations cannot cause infinite execution.

Classes

* Different classes expose different actions.
* Each class receives its configured action-point budget.
* An action unavailable to a class cannot be executed by that class.

YAML validation

* Invalid action references are reported.
* Invalid costs and effect values are rejected.
* Unknown effect types are reported clearly.
* Existing valid game data continues to load.

⸻

14. Implementation Sequence

Complete the implementation as a cohesive task rather than stopping for confirmation after each stage.

1. Analyze the existing repository and identify reusable components.
2. Define the YAML schema and corresponding Go data structures.
3. Implement loading and validation.
4. Implement runtime combat state and action-point accounting.
5. Implement generic action execution and effects.
6. Implement response references and response-chain handling.
7. Implement critical effects and free uses.
8. Integrate class-specific actions.
9. Integrate the CLI.
10. Add and run unit and integration tests.
11. Update relevant documentation and example YAML files.
12. Run the project’s existing test, formatting, and validation commands.

Adjust the sequence when repository constraints require it, while preserving the architectural responsibilities.

Do not replace existing working systems without a concrete technical reason.

⸻

15. Completion Criteria

The implementation is complete when:

* Combat definitions can be loaded from YAML.
* Classes can define different action-point budgets and available actions.
* Actions can consume shared action points and individual uses.
* Actions can trigger other actions through YAML references.
* Critical effects can grant free additional uses.
* The CLI can execute a complete combat exchange.
* Invalid configurations are reported clearly.
* The tests pass.
* Existing game functionality remains intact.

At completion, report:

1. What was implemented.
2. Which files were created or modified.
3. The YAML format and any deviations from the proposal.
4. Tests and validation commands executed, including their results.
5. Any remaining limitations or decisions that could not be resolved from the existing code and requirements.

Do not claim that tests or validation passed unless they were actually executed.