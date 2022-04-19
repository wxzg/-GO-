package main

func main(){

	serve := CreateNewServer("127.0.0.1", 8888)
	serve.Start()
}
