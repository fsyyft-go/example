package testing

func ExamplePrintln() {
	Println(1)
	Println(1, 2)
	// Output:
	// -=>       1
	// -=>       1 2
}

func ExamplePrintf() {
	Printf("abc\n")
	Printf("%[1]s\n", "a")
	Printf("%[1]s%[1]s\n", "a")
	// Output:
	// -=>       abc
	// -=>       a
	// -=>       aa
}
