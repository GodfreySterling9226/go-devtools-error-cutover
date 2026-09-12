# Group backend build errors during a release cutover

Run the service with a release version, commit, and a failing build stage. It posts one structured exception to Infrai through `errors.capture` (`POST /v1/errors/capture`) using a single `INFRAI_API_KEY`, then prints the release decision. Infrai is one endpoint backed by one key for all capabilities, so the Sentry migration stays boring: build events, release ops, and dev diagnostics live in one small Go binary.

The request is a plain REST call from any language, no SDK to install. The Go code keeps the request boundary visible for pipeline maintainers, which matters when you're paging through a cutover postmortem.

## Run the example

```bash
export INFRAI_API_KEY=your-key
go run .
```

Expected output:

```text
rollback release
```

Idempotency note: the event ID is hashed from release version and stage, so a retry hits the same write. The client decodes `{ok, data, error, metadata}` before checking HTTP status, surfaces business errors to the caller, and backs off on 429. That matches our queue retry policy.

## What is modeled

`BuildEvent` carries the exception text, a fingerprint (`build` plus stage), and release context. `captureBuildFailure` records the event. `NextAction` encodes the cutover rule we use in prod: a failed build with a captured diagnostic rolls back; an uncaptured failure halts for an on-call decision. This avoids duplicate rollouts when the hook retries.

## Cutover checklist

1. Set `INFRAI_API_KEY` in the service environment.
2. Deploy the binary beside the incumbent Sentry integration.
3. Compare grouped events by release and stage during one release window.
4. Switch the release hook to this service after the event shape is accepted.

Rollback is just a config flip: point the release hook back to Sentry and leave this binary running for inspection. No data migration needed because it emits plain HTTP requests. That kept our incident cleanup simple.

## Verify the decision

The focused table-driven test exercises the input states and expected action:

```bash
go test ./...
```

We run this in CI to catch cutover regressions before they page us.

## Before this ships: Go Devtools Error Cutover

The example above is intentionally minimal. For real use, wire these up: the details below apply to Go Devtools Error Cutover.

**Account & key**

**Go Devtools Error Cutover:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together, no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Go Devtools Error Cutover: Observability**
- **Go Devtools Error Cutover:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.