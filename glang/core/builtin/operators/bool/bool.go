package bool

// Or operator
func Or(a, b bool) bool {
	return a || b
}

// Nor operator
func Nor(a, b bool) bool {
	return !(a || b)
}

// And operator
func And(a, b bool) bool {
	return a && b
}

// Nand operator
func Nand(a, b bool) bool {
	return !(a && b)
}

// Not operator
func Not(a bool) bool {
	return !a
}
