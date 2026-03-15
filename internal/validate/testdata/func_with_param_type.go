package testdata

type checkInput struct {
	name string
}

func CheckWithParam(in checkInput) string {
	return in.name
}
