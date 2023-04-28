package weave

import "log"

var LogWarning = func(e any) {
	log.Println("Warning:", e) // for now good enough
}
