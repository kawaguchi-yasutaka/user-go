package unixtime

import "time"

type UnixTime int64

func NewUnixTime(t time.Time) UnixTime {
	return UnixTime(t.Unix())
}

func Now() UnixTime {
	return UnixTime(time.Now().Unix())
}
