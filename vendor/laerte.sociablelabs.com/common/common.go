// Package common provides helper methods that can be shared
// throughout the application.
package common

import ()

// error handling
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

/*

*/
func ClearByte(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}
