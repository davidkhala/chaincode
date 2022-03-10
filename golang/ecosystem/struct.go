package main

import . "github.com/davidkhala/goutils"

func panicEcosystem(Type, message string) {
	PanicString("ECOSYSTEM|" + Type + "|" + message)
}
