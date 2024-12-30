package main

func main(){
	app := App{}
	app.initialize()
	app.connect("localhost:10000")
}
