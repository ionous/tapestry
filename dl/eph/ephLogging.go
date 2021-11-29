package eph

import "log"

var LogWarning = func(e error) {
	log.Println("Warning:", e) // for now good enough
}
