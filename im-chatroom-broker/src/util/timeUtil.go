package util

import "time"

func CurrentMillionSecond() int64{
	return time.Now().UnixNano() / 1e6
}

func CurrentSecond() int64{
	return time.Now().Unix()
}
