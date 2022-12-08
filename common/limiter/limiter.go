// Package limiter is to control the links that go into the dispather
package limiter

import (
	//"golang.org/x/time/rate"
	"sync"
)

type Inbound struct {
	UserOnlineIPs sync.Map // Key: Email, Value: [*sync.Map: (Key: IP, Value: UID)]
}

type Limiter struct {
	Inbound Inbound
}

var limiter Limiter

func CheckDeviceLimit(uid uint32, email string, deviceLimit uint32, ip string) bool {
	// Local device limit
	ipMap := new(sync.Map)
	ipMap.Store(ip, uid)
	// If any devices for this email are online
	if v, ok := limiter.Inbound.UserOnlineIPs.LoadOrStore(email, ipMap); ok {
		// Get all current online ip:uid maps for this email
		ipMap := v.(*sync.Map)
		// If this is a new IP
		if _, ok := ipMap.LoadOrStore(ip, uid); !ok {
			// Get the number of online IPs including this new IP
			var counter uint32 = 0
			ipMap.Range(func(key, value interface{}) bool {
				counter++
				return true
			})
			// Delete this new IP if online IPs exceeds the device limit
			if counter > deviceLimit {
				ipMap.Delete(ip)
				return true
			}
		}
	}

	return false
}

func resetDeviceLimit() error {
	limiter.Inbound.UserOnlineIPs.Range(func(key, value interface{}) bool {
		email := key.(string)
		limiter.Inbound.UserOnlineIPs.Delete(email)
		return true
	})

	return nil
}
