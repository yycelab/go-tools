package assert

// !condition ,panic(msg)
func AssertTrue(msg string, condition bool) {
	if !condition {
		panic(msg)
	}
}

// obj == nil ,panic(msg)
func AssertNotNil(msg string, obj any) {
	if obj == nil {
		panic(msg)
	}
}
