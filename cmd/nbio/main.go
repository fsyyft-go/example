package main

func main() {
	if err := rootCmd.Execute(); nil != err {
		panic(err)
	}
}
