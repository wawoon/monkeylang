package object

func MakeInt(i int64) *Integer {
	return &Integer{Value: i}
}
