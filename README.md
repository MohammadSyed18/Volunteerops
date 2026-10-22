# Volunteerops

Volunteer shifts in Go. Handles availability, skills, capacity, overlapping assignments, cancellations, and waitlists, with a small HTTP API.

```sh
go test -race ./...
go run ./cmd/server
```

There are example requests for the three POST endpoints in `sample/`. The data is invented and state stays in memory. Authentication, persistent storage, and organization-specific eligibility rules are still missing.
