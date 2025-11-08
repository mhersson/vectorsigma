# Incremental Updates in VectorSigma

When regenerating a finite state machine with VectorSigma, the tool
intelligently handles existing files to preserve your custom implementations
while updating the necessary components. Here's how VectorSigma treats different
files during regeneration:

## File Treatment During Regeneration

### Files Checked for Changes

These files are analyzed for changes between your implementation and the newly
generated stubs:

- `actions.go`
- `actions_test.go`
- `guards.go`
- `guards_test.go`

**Important:** VectorSigma only processes functions that have the special
comment tag `// +vectorsigma:action:` or `// +vectorsigma:guard:` prefixes. Any
custom helper functions you add to these files without these tags will be left
completely untouched during regeneration. This allows you to safely add your own
utility functions alongside the generated code.

### Files Always Overwritten

These files are completely regenerated each time:

- `zz_generated_statemachine.go` - The core state machine implementation
- `zz_generated_statemachine_test.go` - Help functions for integration testing
  the state machine

### Files Skipped if They Exist

- `extendedstate.go` - Your custom extended state definitions are preserved.
  This file contains the `Context` struct for dependencies like loggers and
  client libraries, and the `ExtendedState` struct for storing data that actions
  modify and guards read. VectorSigma generates this file once but never
  overwrites it, allowing you to add custom fields for your specific
  application's needs without losing them during regeneration.
- `common_test.go` - Common variables and setup code for your tests are
  preserved. This file is generated once to provide a starting point for your
  tests, and can be used by both unit and integration tests. You can add your
  own shared test setup logic here.
- `statemachine_integration_test.go` - Your integration tests are preserved.
  This file is generated once to provide a starting point for your tests. You
  can add your own test cases here to validate the behavior of your state
  machine. Initailly, it only conains a basic happy-path test, and scaffolding
  for optional setup and teardown hooks.

## How Incremental Updates Work

During the regeneration process, VectorSigma follows these steps:

1. Parses your UML diagram to understand the state machine structure
2. Checks for existing implementation files in the output directory
3. For `actions.go` and `guards.go`:
   - Identifies functions with the `// +vectorsigma` comment tag
   - Preserves your existing implementations of tagged methods
   - Adds new methods that weren't present in the previous version
   - Removes methods that are no longer referenced in the UML
   - Never modifies or removes functions without the vectorsigma tag
4. Always regenerates the core state machine files with the `zz_generated_`
   prefix
5. Skips generating `extendedstate.go` if it already exists to preserve your
   custom state

This approach allows you to incrementally update your state machine as your
requirements evolve, without losing your custom implementations or helper
functions.

## Summary

When working with VectorSigma's incremental update feature:

- Do not modify the method signatures in `actions.go` and `guards.go` as they
  must match the generated code exactly
- Add custom helper functions in the actions and guards files as needed
- VectorSigma will preserve your existing implementations of tagged methods
- VectorSigma will never modify or remove functions without the vectorsigma tag
- Use the extended state to share data between actions and guards
- The `extendedstate.go` file is generated once but never overwritten, allowing
  you to add custom fields for your application's needs
- Always review added or removed methods after regeneration to ensure your state
  machine logic remains consistent
