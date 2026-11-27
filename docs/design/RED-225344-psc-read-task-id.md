# PSC read task ID validation

**TLDR:** Reject missing or empty task IDs on PSC read responses before the SDK dereferences them; retain the existing handling of valid IDs and transport errors. This does not address ambiguous outcomes of writes.

Date: 2026-11-27. Status: proposed. Tracking: [RED-225344](https://redislabs.atlassian.net/browse/RED-225344).

## Problem and scope

The PSC service, endpoints, creation-script and deletion-script read methods dereference the task ID returned by the initial GET. A successful response with an absent ID causes a nil-pointer panic. This was reproduced with a simulated response, not observed from the deployed Cloud API. This proposal guards those four read implementations, including their Active-Active entry points. It does not modify create, update, delete or their recovery semantics.

## Behavior and tradeoffs

When a read response has a nil or empty task ID, return a fixed error without the response body and do not invoke the task waiter. Empty strings previously reached the waiter; rejecting them is a compatibility change. Existing transport errors, valid task IDs and the existing 404 translation remain unchanged. Retrying reads has operational costs: Crossplane currently retries unavailable lookups at a fixed one-minute interval. Its prior observation is cleared on error; consumers must handle unavailable status. No new account credentials, permissions, API fields or durable state are introduced.

This is not a safety fix for a mutation whose response is missing: the Cloud operation might already have been accepted. Durable identity, replay prevention, task correlation and shared-account fencing remain separate open issues. It is also not a claim that deployed Cloud returns malformed responses.

## Validation

Hosted [Crossplane experiment #154](https://github.com/RedisLabs/crossplane-provider-rediscloud/pull/154) reproduces the missing-ID panic on SDK main `2cb6b15f705ae2969253d6e3da7d870e7d22cc44` and checks the candidate with a stub waiter, dummy-authenticated loopback HTTP, and actual Crossplane lookup reconcilers with the temporary panic boundary disabled in the disposable CI checkout. The runner replaces Crossplane's pinned SDK with the locally patched SDK main, so it is not an adoption or release compatibility test. Evidence: [run 45990026127](https://github.com/RedisLabs/crossplane-provider-rediscloud/actions/runs/45990026127), including 32 named stub cases, 80 HTTP cases across eight read entry points, and 27 reconciler cases. These cover simulated responses; they do not prove all SDK consumers, cancellation, real Cloud behavior or live-account recovery.

## Rollout and rollback

Review and merge through the SDK's normal process; then release a version and verify dependency compatibility before adopting it in Crossplane. Do not remove the Crossplane containment until the released SDK and relevant read paths have been verified. Roll back the SDK version if valid task IDs or existing error behavior regress. Live acceptance remains gated on account-wide serialization and safe cleanup.

## Open questions

Should a public typed error be exposed to SDK consumers, and is an empty task ID invalid under the Cloud contract? Empty strings previously reached the task waiter; the prior external response is not established. What response and timeout semantics should callers use when read failures persist? Similar task-ID dereferences exist in PSC write paths and other services; they are outside this read-only fix. SDK mutation replay and lost-response recovery require a separate design.
