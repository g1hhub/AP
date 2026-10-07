# Analytics Platform — Product Requirements

## 1. Overview

Analytics Platform v0.1 is an API-first B2B product analytics system. Client applications submit user events, and technical users query aggregated event counts for projects they can access.

The product provides a complete vertical flow:

1. Create a project.
2. Create a user API key.
3. Submit events to the project.
4. Query aggregated analytics for the project.

The initial operator is a technical user who may act as product manager, analyst, architect, and developer. The API is also intended to support backend developers integrating event ingestion and technically proficient analysts querying results. v0.1 has no product-owned frontend; its operations must be described by machine-readable, interactive API documentation.

**Product statement:** Analytics Platform v0.1 enables technical product teams to register events from their applications and retrieve event counts for selected periods through an API, providing a complete path from an isolated project to an analytics result.

## 2. Goals

- Deliver a usable end-to-end analytics workflow rather than ingestion alone.
- Accept and retain application events under a clearly defined contract.
- Provide simple event-count analytics by event name and time period.
- Isolate events and analytics between projects.
- Define explicit validation, duplicate, batch, acceptance, and failure semantics.
- Exercise serious backend requirements through reliability, scale, latency, isolation, security, and observability targets without expanding the product feature set unnecessarily.
- Preserve an event contract that can support more advanced analytics in later versions.

## 3. Non-Goals

- Providing a graphical analytics dashboard.
- Reproducing the full feature set of Amplitude or Mixpanel.
- Supporting non-technical end users in v0.1.
- Defining storage, messaging, deployment, orchestration, or monitoring technologies.
- Providing self-service user registration, interactive login, or profile management.
- Implementing organization-level tenancy in v0.1.
- Supporting advanced analytics or property-based analysis in v0.1.
- Providing commercial service-level agreements; scale and latency values in this document are educational design targets.

## 4. Reference Use Case

`CoolShop` is an online store whose application emits events such as:

```text
signup
  -> search
  -> view_product
  -> add_to_cart
  -> checkout_started
  -> purchase
```

CoolShop creates separate projects when it needs isolated datasets, for example `CoolShop Production` and `CoolShop Staging`. Its application submits events to the appropriate project. A user with access to that project can then ask:

- How many `purchase` events occurred during a selected period?
- How many events of each event type occurred during a selected period?

CoolShop is only a reference scenario. Projects, events, and properties must not contain e-commerce-specific semantics.

## 5. Domain Model

### User

A User represents a human or operator identity that accesses the API.

- v0.1 supports multiple users.
- Users are provisioned by the system operator outside the product API.
- A user may access multiple projects.
- Multiple users may access the same project.
- A user authenticates API requests with a user-owned API key.
- The event field `user_id` identifies a user of the client application and is unrelated to this platform User.

### Organization

Organization is not part of the v0.1 domain model. Project is the highest tenancy and data-isolation boundary. A future version may introduce organizations to represent customer companies and group users and projects.

### Project

A Project is a universal, isolated data scope. A client decides whether a project represents an application, environment, team, or another separation boundary.

- A user may create and view projects through the API.
- The creating user automatically receives access to the project with the informational `owner` role.
- A project may have multiple users, and a user may have multiple projects.
- Assigning additional users to projects is performed by the system operator outside the product API.
- Project memberships carry `owner` or `member` roles. These roles are informational in v0.1 and grant identical permissions.
- Every event belongs to exactly one project.
- Analytics is evaluated within exactly one project.
- Cross-project ingestion and analytics are not supported.
- Separate projects are used when production and staging data must be isolated.
- v0.1 supports project creation and viewing, but not project update or deletion.

### API Key

An API Key is a user-owned credential and is the only API authentication mechanism in v0.1.

- A key grants its user all v0.1 operations, including project creation, project viewing, event ingestion, analytics queries, and management of the user's own keys.
- A key applies to every project accessible to its user; it is not scoped to one project.
- The project is identified separately for project-specific operations, and the user's access to it is checked.
- A user may have multiple simultaneously active keys.
- An authenticated user may create and revoke only their own keys.
- Key creation returns a `key_id` and reveals the secret exactly once.
- Existing secrets cannot be retrieved again.
- v0.1 does not provide an operation to list keys. A user must retain `key_id` to revoke a key.
- Revocation is irreversible and takes effect immediately for all subsequent requests.
- A key may revoke itself or the user's last active key. The operator must restore access if the user loses all active credentials.
- Initial credential bootstrap is performed by the operator outside v0.1 requirements.
- The intentionally accepted v0.1 trade-off is that compromise of one key exposes every operation and project available to its user.

### Event

An Event records an occurrence reported by a client application within one project.

An event contains:

- a required normalized event name;
- a required client-assigned `user_id`;
- a system-generated `event_id`;
- an optional client-assigned event time;
- a system-assigned receipt time;
- optional JSON properties.

Events are immutable in v0.1. They cannot be updated or deleted through the product API.

## 6. Event Contract

### Required Fields

#### `event`

- Purpose: identifies the type of occurrence and is the grouping key for v0.1 analytics.
- Supplied by the client.
- Must be a Unicode string.
- Leading and trailing whitespace is removed, then the value is converted to Unicode lowercase.
- The normalized value must not be empty.
- Internal whitespace is preserved. For example, `product viewed` and `product   viewed` are different names.
- The normalized value is stored and used for aggregation. `Purchase`, `PURCHASE`, and `" purchase "` therefore identify the same event type, `purchase`.
- Other than trimming, lowercasing, non-emptiness, and the length limit, v0.1 imposes no character restriction.

#### `user_id`

- Purpose: identifies the user in the client application that produced the event.
- Supplied by the client and required for every event.
- Treated as an opaque string; the platform does not interpret its format.
- Its identity scope is one project. Equal values in different projects do not imply the same person.
- Anonymous identity, `anonymous_id`, and `device_id` are not supported in v0.1.

### Optional Fields

#### Client event time

The client may provide event time using exactly one of these fields:

- `event_time_seconds`: an integer Unix timestamp in seconds; or
- `event_time_milliseconds`: an integer Unix timestamp in milliseconds.

If both fields are supplied, or a supplied value is invalid, the event is rejected. If neither is supplied, the event remains valid and its system-assigned `received_at` is used as its effective event time.

#### `properties`

`properties` is an optional JSON object containing client-defined event attributes. When absent, the event has no properties.

### System-Assigned Fields

#### `event_id`

- Generated by the platform for every accepted event.
- Uniquely identifies the stored event within the platform.
- Returned to the client after successful acceptance.
- Cannot be supplied or selected by the client and therefore cannot be used by the client as an idempotency key.

#### `received_at`

- Generated by the platform for every accepted event.
- Records when the platform received the event.
- Cannot be supplied or overridden by the client.

### Properties

- The root value of `properties` must be a JSON object.
- Supported values are strings, JSON numbers, booleans, `null`, arrays, and nested objects.
- Arrays may contain arbitrary supported JSON values.
- All valid JSON numbers are accepted. Non-JSON values such as `NaN` and positive or negative infinity are rejected.
- Nested arrays and objects are limited by the maximum nesting depth.
- v0.1 accepts and retains properties for future use but does not filter, group, or aggregate by them.
- Except for the defined length and size limits, v0.1 places no character restrictions on property keys.

### Validation Rules

An event is valid only when:

- all required fields are present and have the required types;
- the normalized event name is non-empty and within its length limit;
- `user_id` is present and is a string;
- at most one client event-time field is supplied and its value is a valid integer Unix timestamp in the field's declared unit;
- `properties`, when supplied, is a JSON object containing only supported JSON values;
- all property count, key length, string size, nesting depth, event size, batch size, and request size limits are satisfied.

An invalid single event is rejected. A batch containing any invalid event is rejected in full; no event from that batch is accepted.

### Timestamp Semantics

- `received_at` is always assigned by the platform.
- Effective event time is the client event time when supplied; otherwise it is `received_at`.
- Analytics periods are evaluated against effective event time.
- Valid client event times are accepted regardless of how far they are in the past or future.
- Consequently, an erroneous but valid future timestamp is counted in that future analytics period; the platform does not silently replace it.
- Analytics periods are half-open intervals `[from, to)`, including `from` and excluding `to`.
- Analytics query boundaries are Unix timestamps with an explicitly identified unit. The API design must not infer seconds versus milliseconds from numeric magnitude.
- A query with `from >= to` is rejected.

### Duplicate Semantics

- Duplicate client submissions are allowed and are counted as separate events.
- Every accepted submission receives a new system-generated `event_id`.
- If a client times out and retries the same event, both attempts may be accepted and counted.
- v0.1 provides neither best-effort deduplication nor strict idempotency.
- Clients must not interpret matching content as event identity.

## 7. Core User Flows

### Project Creation

1. A pre-provisioned user authenticates with an active API key.
2. The user submits a project creation request.
3. The platform creates an isolated project.
4. The creating user receives project access with the informational `owner` role.
5. The user can view the created project and use it in project-specific requests.

Additional user assignment is an operator action outside the API.

### API Key Creation

1. A user authenticates with an existing active API key, except for operator-managed initial bootstrap.
2. The user requests creation of a key for themselves.
3. The platform returns `key_id` and displays the secret once.
4. The user retains `key_id` if future revocation is required.
5. The user may later revoke the key by `key_id`, including the current or last active key.

### Event Ingestion

1. The client authenticates with a user API key and identifies the target project.
2. The platform verifies that the key is active and its user has access to the project.
3. The client submits one event or an atomic batch of events.
4. The platform validates authentication, project access, every event, and all request limits.
5. If any validation fails, the request is rejected according to the error contract. No event from an invalid batch is accepted.
6. On success, the platform accepts all events and returns generated event IDs in the same order as the input events.

### Analytics Query

1. A user authenticates with a user API key and identifies a project they can access.
2. The user supplies a valid half-open time period `[from, to)` with explicit timestamp units.
3. The user requests either the count for one normalized event type or counts grouped by event type.
4. The platform evaluates effective event times only within the selected project.
5. No matching event produces a successful zero count or empty grouped result, not an error.

## 8. Functional Requirements

### FR-001 Project creation

The system shall allow an authenticated user to create a project.

### FR-002 Creator membership

The system shall grant the project creator access to the created project with the informational `owner` role.

### FR-003 Project viewing

The system shall allow a user to view only projects to which the user has been assigned.

### FR-004 Project membership

The system shall support multiple users per project and multiple projects per user.

### FR-005 Informational roles

The system shall retain `owner` and `member` project membership roles without applying different permissions to them in v0.1.

### FR-006 User key creation

The system shall allow an authenticated user to create multiple API keys for themselves.

### FR-007 One-time secret disclosure

The system shall return a new key's `key_id` and secret at creation and shall not expose the secret again.

### FR-008 User key revocation

The system shall allow a user to revoke only their own API key by `key_id`, including the currently used or last active key.

### FR-009 Immediate revocation

The system shall reject every request made with a key after that key has been revoked.

### FR-010 Universal user-key access

An active API key shall authorize all v0.1 operations and projects available to its owning user.

### FR-011 Project access verification

The system shall verify the key owner's project membership for every project-specific operation.

### FR-012 Single-event ingestion

The system shall accept a valid single event for a project accessible to the authenticated user.

### FR-013 Batch ingestion

The system shall accept a valid batch of events for one project up to the documented limits.

### FR-014 Atomic batch validation

The system shall reject an entire batch when any event in it is invalid or any request limit is exceeded.

### FR-015 Event identity

The system shall generate and return a distinct `event_id` for every accepted event.

### FR-016 Receipt time

The system shall record a system-assigned `received_at` for every accepted event.

### FR-017 Event-time fallback

The system shall use client event time for analytics when supplied and `received_at` otherwise.

### FR-018 Event-name normalization

The system shall trim leading and trailing whitespace and convert event names to lowercase before validation, retention, and aggregation.

### FR-019 Property retention

The system shall accept and retain supported event properties without using them in v0.1 analytics.

### FR-020 Duplicate acceptance

The system shall treat repeated submissions as independent events and assign each a distinct `event_id`.

### FR-021 Count by event type

The system shall return the number of events of a selected normalized event type within a project and half-open time period.

### FR-022 Counts grouped by event type

The system shall return event counts grouped by normalized event type within a project and half-open time period.

### FR-023 Empty analytics result

The system shall return a successful zero or empty grouped result when no events match a valid analytics query.

### FR-024 Error distinction

The system shall return errors that allow clients to distinguish authentication, project-access, validation, and internal failures.

### FR-025 API documentation

The system shall provide machine-readable, interactive documentation for all v0.1 API operations and contracts.

## 9. Non-Functional Requirements

### NFR-001 Ingestion latency

At target load, successful ingestion requests shall complete within 200 ms at the 95th percentile.

### NFR-002 Analytics latency

At target scale, supported analytics queries shall complete within 2 seconds at the 95th percentile.

### NFR-003 Analytics freshness

Under normal operation, an accepted event shall become visible to analytics queries within 60 seconds.

### NFR-004 Accepted-event durability

After a successful ingestion response, an accepted event shall survive ordinary failure of an individual application process.

### NFR-005 Honest acknowledgement

The system shall not report an ingestion request as successful unless all events covered by that response satisfy the acceptance semantics in this document.

### NFR-006 Horizontal growth

The system design shall permit growth beyond the target throughput without changing the public event contract.

### NFR-007 Uneven tenant load

The system shall tolerate a project generating a significant share of total traffic without weakening project isolation.

### NFR-008 Tenant confidentiality

The system shall prevent users from reading or writing project data unless they have access to that project.

### NFR-009 Credential confidentiality

API key secrets shall not be exposed after initial creation or included in logs, errors, metrics, traces, or other diagnostic output.

### NFR-010 Operational visibility

The system shall expose operational information sufficient to detect ingestion volume changes, failures, increased latency, rejected events, and delayed analytics freshness.

### NFR-011 Diagnosability

Operational information and client errors shall be sufficient to distinguish authentication, access, validation, and internal failure classes.

### NFR-012 Maintainability

Ingestion and analytics capabilities shall be evolvable independently without requiring a breaking change to the public event contract.

### NFR-013 Contract documentation

Validation rules, limits, response semantics, and errors shall be documented consistently with observable API behavior.

## 10. Limits

The following v0.1 limits make the input contract finite and reject pathological requests. They are product limits rather than an implementation prescription.

| Item | Limit | Rationale |
| --- | ---: | --- |
| Normalized event name | 128 Unicode characters | Allows descriptive names while preventing unbounded grouping keys. |
| Property key | 128 Unicode characters | Allows descriptive metadata without unbounded keys. |
| Top-level properties | 100 | Supports useful event context while bounding per-event complexity. |
| Nested arrays and objects | 5 container levels | Supports structured metadata while preventing pathologically deep documents. |
| Individual string property value | 8 KiB encoded size | Allows normal descriptive values while discouraging document or binary payloads. |
| Encoded event | 64 KiB | Bounds validation, transfer, retention, and processing cost per event. |
| Events per batch | 500 | Provides useful batching while bounding atomic validation and response size. |
| Encoded ingestion request | 5 MiB | Bounds total work and memory associated with one request. |

- An event or request exceeding any applicable limit is invalid.
- The whole batch is rejected if one event or the request exceeds a limit.
- Overall event and request size limits apply in addition to all field-specific limits.
- An empty batch is invalid because it contains no event to accept.

## 11. Tenant Isolation

- Project is the highest tenant and data-isolation boundary in v0.1.
- Every accepted event belongs to exactly one explicitly selected project.
- A key identifies its owning user, not a project. The target project must be identified separately.
- The key owner must have project membership for ingestion, project viewing, and analytics operations.
- A user with access to Project A but not Project B cannot submit events to or read analytics from Project B.
- Data from Project A must never contribute to a query for Project B.
- A request cannot submit one batch to multiple projects.
- Cross-project analytics is not supported.
- Membership roles do not weaken these rules; `owner` and `member` have the same project data access in v0.1.

## 12. Acceptance Semantics

A successful ingestion response means:

- authentication and project access checks succeeded;
- the single event or every event in the batch passed validation;
- all events in the request were accepted as one complete outcome;
- generated event IDs are returned in input order;
- accepted events shall not be lost due to ordinary failure of an individual application process.

A successful response does not mean:

- the events are already visible in analytics;
- a retried event will be deduplicated;
- events are protected from every catastrophic or correlated failure;
- future timestamps have been corrected or old timestamps have been rejected.

In normal operation, accepted events must become queryable within the 60-second analytics freshness target.

## 13. MVP Analytics

v0.1 supports exactly two analytics operations within one project:

1. Count events of one normalized event type over `[from, to)`.
2. Count all matching events over `[from, to)`, grouped by normalized event type.

Analytics uses effective event time, defined as client event time when supplied and `received_at` otherwise.

v0.1 analytics does not:

- count unique users;
- calculate DAU, WAU, or MAU;
- filter or group by properties;
- calculate funnels, retention, or cohorts;
- combine projects;
- promise immediate visibility after ingestion.

## 14. Target Scale

These values are educational design targets for future architecture work, not claims derived from real customer traffic and not external SLAs.

| Dimension | Target |
| --- | ---: |
| Average ingestion throughput | 1,000 events/second |
| Peak ingestion throughput | 10,000 events/second |
| Projects | 1,000 |
| Accumulated events across all projects | 1 billion |
| Ingestion latency at target load | p95 <= 200 ms |
| MVP analytics query latency at target scale | p95 <= 2 seconds |
| Analytics freshness in normal operation | <= 60 seconds |

Load may be uneven across projects. Future scaling beyond these targets should not require a breaking change to the event contract.

## 15. Out of Scope

The following are explicitly excluded from v0.1:

- frontend dashboard;
- self-service registration, interactive login, and user profile management;
- Organization;
- user provisioning and project membership management through the product API;
- project update and deletion;
- removal of project members through the product API;
- permission differences between `owner` and `member`;
- scoped key types or restricted key permissions;
- API key listing and repeated secret disclosure;
- SDKs;
- property filters and breakdowns;
- unique users, DAU, WAU, and MAU;
- funnels, retention, and cohorts;
- cross-project analytics;
- event deduplication and client-provided idempotency keys;
- event update and deletion;
- real-time streaming analytics;
- billing, quotas, and pricing plans;
- alerting;
- SSO;
- data export;
- GDPR-specific tooling;
- arbitrary SQL access;
- guaranteed event processing order;
- guarantees against all catastrophic or correlated failures.

## 16. Open Questions

The following questions are intentionally deferred beyond v0.1:

- Should Organization be introduced to group customer users and projects?
- Should users be able to register, authenticate interactively, and manage project membership without an operator?
- Should keys be separated or scoped by purpose and permission?
- Should anonymous identity and later association with a known user be supported?
- Should the client provide an idempotency identifier for retry deduplication?
- Which property filters and breakdowns should be introduced first?
- What event retention and deletion policy should apply?
- Are per-user or per-project ingestion quotas required?
- When should `owner` and `member` receive distinct permissions?
- Should erroneous future or excessively old events be correctable or rejectable?
- What disaster-recovery guarantees are required for catastrophic failures?
- Should Unicode event names receive canonical normalization beyond trimming and lowercasing?
