# volunteerops

volunteer shifts in Go. handles availability, skills, capacity, overlapping assignments, cancellations, and waitlists, with a small HTTP API.

```sh
go test -race ./...
go run ./cmd/server
```

there are example requests for the three post endpoints in `sample/`. the data is invented and state stays in memory. authentication, persistent storage, and organization-specific eligibility rules are still missing.

## Implementation details

Waitlist promotions recheck current availability and assignment conflicts. Individual waitlist requests can be withdrawn, and cancelling a shift releases all of its assignments exactly once.
