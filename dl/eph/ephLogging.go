package eph

import "log"

func LogWarning(e error) {
	log.Println("Warning:", e) // for now good enough
}
