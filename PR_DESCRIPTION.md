### What problem does this PR solve?

Issue Number: Close #2810

Problem Summary:
Users have requested an example of how to integrate `json-iterator` with Gin at runtime (without using build tags). While there is documentation, a runnable example project is helpful for understanding the integration points, specifically implementing the `json.Core` interface and replacing `json.API`.

### What is changed and how does it work?

- Added a new example under `examples/json-iterator`.
- Implemented `customJsonApi` which wraps `jsoniter.Config` and implements `json.Core`.
- Demonstrates how to replace the default `json.API` with the custom implementation.
- Added a unit test to verify the integration works as expected.
- Added a README for the example explaining how to run and test it.

### Check List

Tests

- [x] Unit test
  - Added `examples/json-iterator/json_iterator_test.go`
- [ ] Integration test
- [x] Manual test
  - Verified with `curl` locally.
- [ ] No code

Code changes

- [ ] Has the configuration change
- [ ] Has HTTP API interfaces changed
- [ ] Has persistent data change

Side effects

- [ ] Possible performance regression
- [ ] Increased code complexity
- [ ] Breaking backward compatibility

Related changes

- [ ] PR to update [`pingcap/docs`](https://github.com/pingcap/docs)/[`pingcap/docs-cn`](https://github.com/pingcap/docs-cn):
- [ ] PR to update [`pingcap/tiup`](https://github.com/pingcap/tiup):
- [ ] Need to cherry-pick to the release branch

### How to Test

1. Navigate to the example directory:
   ```bash
   cd examples/json-iterator
   ```
2. Run the tests:
   ```bash
   go test -v
   ```
3. Run the example:
   ```bash
   go run main.go
   ```
4. Make a request:
   ```bash
   curl http://localhost:8080/ping
   ```
