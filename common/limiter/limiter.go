// Package limiter is to control the links that go into the dispather
package limiter

import (
	//"golang.org/x/time/rate"
	"sync"
)

type InboundInfo struct {
	Tag          string
	BucketHub    *sync.Map // Key: Email, Value: *rate.Limiter
	UserOnlineIP *sync.Map // Key: Email, Value: [*sync.Map: (Key: IP, Value: UID)]
}

type Limiter struct {
	InboundInfo *sync.Map // Key: Tag, Value: *InboundInfo
}

var limiter = &Limiter{
	InboundInfo: new(sync.Map),
}

func CheckDeviceLimit(uid int, email string, deviceLimit int, ip string) bool {
		// Local device limit
		ipMap := new(sync.Map)
		ipMap.Store(ip, uid)
		// If any devices for this email are online
		if v, ok := inboundInfo.UserOnlineIP.LoadOrStore(email, ipMap); ok {
			// Get all current online ip:uid maps for this email
			ipMap := v.(*sync.Map)
			// If this is a new IP
			if _, ok := ipMap.LoadOrStore(ip, uid); !ok {
				// Get the number of online IPs including this new IP
				counter := 0
				ipMap.Range(func(key, value interface{}) bool {
					counter++
					return true
				})
				// Delete this new IP if online IPs exceeds the device limit
				if counter > deviceLimit && deviceLimit > 0 {
					ipMap.Delete(ip)
					return true
				}
			}
		}
}

func resetDeviceLimit() error {
	inboundInfo.UserOnlineIP.Range(func(key, value interface{}) bool {
		email := key.(string)
		inboundInfo.UserOnlineIP.Delete(email) // Reset online device
		return true
	})

	return nil
}
