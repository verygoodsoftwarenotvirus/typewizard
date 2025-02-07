package a

import (
	"time"

	"typewizard/utils/test_packages/example/c"
)

type (
	ThingOne struct {
		Name        string
		Age         int
		DateOfBirth time.Time
		Else        c.SomethingElse
	}
)
