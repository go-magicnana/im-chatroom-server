package util

func Panic(obj interface{}) {

	if obj != nil {
		panic(obj)
	}
}
