# volunteerops

volunteer shifts in Go. handles availability, skills, capacity, overlapping assignments, cancellations, and waitlists, with a small HTTP API.

```sh
go test -race ./...
go run ./cmd/server
```

there are example requests for the three post endpoints in `sample/`. the data is invented and state stays in memory. the HTTP demo does not add authentication or organization-specific eligibility policies automatically.

## Implementation details

Waitlist promotions recheck current availability and assignment conflicts. Individual waitlist requests can be withdrawn, and cancelling a shift releases all of its assignments exactly once.

## Snapshot persistence

[snapshot.go](snapshot.go) exports a consistent detached state and restores it through validated coordinator operations. Restore checks capacities, required skills, overlapping assignments, waitlist references, duplicate requests, cancelled shifts, and recomputed assignment counts.

[snapshot_file.go](snapshot_file.go) writes snapshots through a flushed temporary file and atomic replacement. Imports are bounded to 16 MiB. The application should serialize saves to a given path; the snapshot operation itself is protected by the coordinator mutex.

```sh
go test -race ./...
go run ./cmd/statecheck path/to/snapshot.json
```

Call `coordinator.Save(path)` and `volunteerops.Load(path)` from a host application. The HTTP demo still keeps its active coordinator in memory and does not automatically save after every request. The state-check command validates a saved snapshot and emits shift coverage reports. Race tests exercise snapshots taken during concurrent requests.
