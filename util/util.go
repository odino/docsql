package util

// The ever-missing function to rule 'em all!
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
