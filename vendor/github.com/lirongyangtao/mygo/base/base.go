package base

const (
	E1LessE2     int32 = iota //e1<e2
	E1EqualE2                 //e1=e2
	E1GenerateE2              //e1>e2
)

type CompareAble interface {
	Less(e1 interface{}, e2 interface{}) int32
}

type CompareFunc func(e1 interface{}, e2 interface{}) int32

func CmInt(v1, v2 interface{}) int32 {
	if v1.(int) < v2.(int) {
		return E1LessE2
	} else if v1.(int) > v2.(int) {
		return E1GenerateE2
	} else {
		return E1EqualE2
	}
}
