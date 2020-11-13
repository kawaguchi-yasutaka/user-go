package timekeeper

import (
	"user-go/domain/interfaces"
	"user-go/lib/unixtime"
)

type TimeKeeperMock struct {
	N unixtime.UnixTime
}

func NewTimeKeeperMock() TimeKeeperMock {
	return TimeKeeperMock{}
}

var _ interfaces.ITimeKeeper = TimeKeeper{}

func (t TimeKeeperMock) Now() unixtime.UnixTime {
	return t.N
}
