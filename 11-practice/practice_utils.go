package practice

func Experiment() {
	tprivate := practice.testPrivate{
		typed: "private",
		count: 1,
	}
	fmt.Println(tprivate)

	tpublic := practice.TestPublic{
		Name: "public",
		Age:  30,
	}
	fmt.Println(tpublic)

	tmix := practice.TestMix{
		Name: "mix",
		age:  25,
	}
	fmt.Println(tmix)
}