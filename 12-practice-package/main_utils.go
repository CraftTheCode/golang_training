package practice

type testPrivate struct {
	typed string
	count int
}

type TestPublic struct {
	Name string
	Age  int
}

type TestMix struct {
	Name string
	age  int
}

func GetTestMix(name string, age int) *TestMix {
	return &TestMix{
		Name: name,
		age:  age,
	}
}
