package timekeeper

import (
	"user-go/domain/interfaces"
	"user-go/lib/unixtime"
)

type TimeKeeper struct {
}

var _ interfaces.ITimeKeeper = TimeKeeper{}

func NewTimeKeeper() TimeKeeper {
	return TimeKeeper{}
}

func (t TimeKeeper) Now() unixtime.UnixTime {
	return unixtime.Now()
}
