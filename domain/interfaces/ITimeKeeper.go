package interfaces

import "user-go/lib/unixtime"

type ITimeKeeper interface {
	Now() unixtime.UnixTime
}
