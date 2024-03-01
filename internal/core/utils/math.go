package utils

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		return *new(T) // zero value of T
	}

	if isNan(args[0]) {
		return args[0]
	}

	max := args[0]
	for _, arg := range args[1:] {

		if isNan(arg) {
			return arg
		}

		if arg > max {
			max = arg
		}
	}
	return max
}

func isNan[T constraints.Ordered](arg T) bool {
	return arg != arg
}
