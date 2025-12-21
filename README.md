# Volunteerops

A Go service for concurrency-safe volunteer and shift coordination. The core enforces availability, skills, capacity, overlap conflicts, cancellations, waitlist promotion, age/training/workload policies, fair candidate ranking, weekly-minute calculations, and shift fill reports. A small HTTP API exposes volunteer, shift, and request creation.

Run `go test -race ./...` and `go run ./cmd/server`. The `sample/` requests contain invented demonstration data for the three POST endpoints. State is intentionally in memory; authentication, durable event storage, and organization-specific eligibility policies remain future work.
