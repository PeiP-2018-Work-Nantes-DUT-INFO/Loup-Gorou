package secondstimer

import "time"

type SecondsTimer struct {
	Timer *time.Timer
	end   time.Time
}

func NewSecondsTimer(t time.Duration) *SecondsTimer {
	return &SecondsTimer{time.NewTimer(t), time.Now().Add(t)}
}

func (s *SecondsTimer) Reset(t time.Duration) {
	s.Timer.Reset(t)
	s.end = time.Now().Add(t)
}

func (s *SecondsTimer) Stop() {
	s.Timer.Stop()
}

func (s *SecondsTimer) TimeRemaining() time.Duration {
	return time.Until(s.end)
}
