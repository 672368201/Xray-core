// Package limiter is to control the links that go into the dispather
package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

type InboundInfo struct {
	Tag          string
	BucketHub    *sync.Map // key: Email, value: *rate.Limiter
	UserOnlineIP *sync.Map // Key: Email Value: *sync.Map: Key: IP, Value: UID
}

type Limiter struct {
	InboundInfo *sync.Map // Key: Tag, Value: *InboundInfo
}

func New() *Limiter {
	return &Limiter{
		InboundInfo: new(sync.Map),
	}
}

func (l *Limiter) GetUserBucket(tag string, email string, deviceLimit int, speedLimit uint64, ip string, port uint16) (limiter *rate.Limiter, SpeedLimit bool, Reject bool) {
	if value, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := value.(*InboundInfo)

		// Local device limit
		ipMap := new(sync.Map)
		ipMap.Store(ip, port)
		// If any device is online
		if v, ok := inboundInfo.UserOnlineIP.LoadOrStore(email, ipMap); ok {
			ipMap := v.(*sync.Map)
			// If this ip is a new device
			if _, ok := ipMap.LoadOrStore(ip, port); !ok {
				counter := 0
				ipMap.Range(func(key, value interface{}) bool {
					counter++
					return true
				})
				if counter > deviceLimit && deviceLimit > 0 {
					ipMap.Delete(ip)
					return nil, false, true
				}
			}
		}

		// If need the Speed limit
		if speedLimit > 0 {
			limiter := rate.NewLimiter(rate.Limit(speedLimit), int(speedLimit)) // Byte/s
			if v, ok := inboundInfo.BucketHub.LoadOrStore(email, limiter); ok {
				bucket := v.(*rate.Limiter)
				return bucket, true, false
			} else {
				return limiter, true, false
			}
		} else {
			return nil, false, false
		}
	} else {
		newError("Get Inbound Limiter information failed").AtDebug().WriteToLog()
		return nil, false, false
	}
}
