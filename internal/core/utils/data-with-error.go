package utils

/*
DataWithError is a generic type for some data that can be returned from a function along with an error.
*/
type DataWithError[T any] struct {
	Data T
	Err  error
}
