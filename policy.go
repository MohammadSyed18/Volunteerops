package volunteerops
import("errors";"sort";"time")
type EligibilityPolicy struct{MinimumAge int;RequiredTraining []string;MaximumWeeklyMinutes int}
type Candidate struct{VolunteerID string;Age int;Training []string;AssignedMinutes int}
func(p EligibilityPolicy)Check(c Candidate)error{if c.Age<p.MinimumAge{return errors.New("minimum age not met")};for _,required:=range p.RequiredTraining{if !contains(c.Training,required){return errors.New("required training missing: "+required)}};if p.MaximumWeeklyMinutes>0&&c.AssignedMinutes>=p.MaximumWeeklyMinutes{return errors.New("weekly assignment limit reached")};return nil}
type Match struct{VolunteerID string;Score int;Reason string}
func RankCandidates(shift Shift,candidates []Candidate,volunteers map[string]Volunteer,policy EligibilityPolicy)[]Match{var out []Match;for _,candidate:=range candidates{v,ok:=volunteers[candidate.VolunteerID];if !ok||policy.Check(candidate)!=nil||v.AvailableFrom.After(shift.Start)||v.AvailableTo.Before(shift.End){continue};skillMatches:=0;for _,s:=range shift.RequiredSkills{if contains(v.Skills,s){skillMatches++}};if skillMatches<len(shift.RequiredSkills){continue};score:=100-candidate.AssignedMinutes/10+skillMatches*10;out=append(out,Match{candidate.VolunteerID,score,"eligible, available, and skill-matched"})};sort.SliceStable(out,func(i,j int)bool{return out[i].Score>out[j].Score});return out}
func WeeklyMinutes(shifts []Shift,id string)int{total:=0;for _,s:=range shifts{if contains(s.Assigned,id){total+=int(s.End.Sub(s.Start)/time.Minute)}};return total}
