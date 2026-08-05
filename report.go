package volunteerops
import("encoding/json";"time")
type ShiftReport struct{ShiftID string `json:"shiftId"`;Filled int `json:"filled"`;Capacity int `json:"capacity"`;Waitlisted int `json:"waitlisted"`;FillRate float64 `json:"fillRate"`;StartsAt time.Time `json:"startsAt"`}
func Summarize(s Shift)ShiftReport{return ShiftReport{s.ID,len(s.Assigned),s.Capacity,len(s.Waitlist),float64(len(s.Assigned))/float64(s.Capacity),s.Start}}
func ReportJSON(s Shift)([]byte,error){return json.MarshalIndent(Summarize(s),"","  ")}
