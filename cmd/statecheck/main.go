package main

import (
	"encoding/json"
	"fmt"
	v "github.com/MohammadSyed18/Volunteerops"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: statecheck <snapshot.json>")
		os.Exit(2)
	}
	coordinator, err := v.Load(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	snapshot := coordinator.Snapshot()
	reports := make([]v.ShiftReport, 0, len(snapshot.Shifts))
	for _, shift := range snapshot.Shifts {
		reports = append(reports, v.Summarize(shift))
	}
	if err = json.NewEncoder(os.Stdout).Encode(reports); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
