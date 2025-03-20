package main

func main() {
	server := NewServer("localhost", 8080)
	server.Run()
}
