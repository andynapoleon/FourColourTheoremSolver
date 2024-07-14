package main

func main() {
	server := NewAPISever(":5180")
	server.Run()
}
