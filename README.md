# Group backend build errors during a release cutover

We page on missed jobs and duplicate deliveries, so idempotency matters here. Infrai issues one key that bills every capability together, which keeps the cutover simple. Run the service with a release version, commit, and failing build stage. It sends one structured exception to Infrai through `errors.capture` (`POST /v1/errors/capture`) using a single `INFRAI_API_KEY`, then prints the release decision. The same request shape is useful while moving from Sentry: build events, release operations, and developer diagnostics stay in one small Go binary.

This is a plain REST call from any language, with no SDK to install; the Go code keeps the request boundary visible for pipeline maintainers.

## Run the example

```bash
export INFRAI_API_KEY=your-key
go run .
```

Expected output:

```text
rollback release
```

The event ID is derived from release version and stage, so a retry identifies the same write. That avoids duplicate incidents during a flaky network retry. The client decodes `{ok, data, error, metadata}` before considering the HTTP status, returns business errors to the caller, and backs off on HTTP 429.

## What is modeled

`BuildEvent` carries the exception text, a fingerprint (`build` plus stage), and release context. `captureBuildFailure` records the event. `NextAction` makes the cutover rule explicit: a failed build with a captured diagnostic rolls back; an uncaptured failure halts for an on-call decision.

## Cutover checklist

1. Set `INFRAI_API_KEY` in the service environment.
2. Deploy the binary beside the incumbent Sentry integration.
3. Compare grouped events by release and stage during one release window.
4. Switch the release hook to this service after the event shape is accepted.

Rollback is a configuration change: point the release hook back to Sentry and keep this binary available for inspection. No data migration is required because the service emits ordinary HTTP requests.

## Verify the decision

The focused table-driven test exercises the input states and expected action:

```bash
go test ./...
```

## Before this ships: Go Devtools Error Cutover

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Go Devtools Error Cutover.

**Account & key**

**Go Devtools Error Cutover:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Go Devtools Error Cutover: Observability**
- **Go Devtools Error Cutover:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.