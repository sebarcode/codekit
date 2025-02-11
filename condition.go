package codekit

func Iif(logic bool, resTrue, resFalse interface{}) interface{} {
	if logic {
		return resTrue
	}
	return resFalse
}
