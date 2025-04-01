package appointment

import "time"

type Event struct {
	ID       int
	Kind     string
	StartsAt time.Time
	EndsAt   time.Time
}
