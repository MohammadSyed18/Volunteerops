package volunteerops
import("sync";"testing";"time")
func fixture()*Coordinator{c:=New();start:=time.Date(2030,1,1,9,0,0,0,time.UTC);c.AddShift(Shift{ID:"s",Start:start,End:start.Add(time.Hour),Capacity:1,RequiredSkills:[]string{"checkin"}});for _,id:=range []string{"a","b","c"}{c.AddVolunteer(Volunteer{ID:id,Skills:[]string{"checkin"},AvailableFrom:start,AvailableTo:start.Add(2*time.Hour)})};return c}
func TestCapacityWaitlistAndPromotion(t *testing.T){c:=fixture();if state,_:=c.Request("a","s");state!="assigned"{t.Fatal(state)};if state,_:=c.Request("b","s");state!="waitlisted"{t.Fatal(state)};if err:=c.Cancel("a","s");err!=nil{t.Fatal(err)};s,_:=c.Shift("s");if len(s.Assigned)!=1||s.Assigned[0]!="b"{t.Fatal(s)}}
func TestConcurrentRequestsNeverExceedCapacity(t *testing.T){c:=fixture();var wg sync.WaitGroup;for _,id:=range []string{"a","b","c"}{wg.Add(1);go func(){defer wg.Done();c.Request(id,"s")}()};wg.Wait();s,_:=c.Shift("s");if len(s.Assigned)!=1||len(s.Waitlist)!=2{t.Fatal(s)}}
