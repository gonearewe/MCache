package main

func main() {
	cache := NewCache()
	NewServer(cache).Listen()
}
