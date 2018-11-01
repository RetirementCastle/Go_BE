// main.go

package main

func main() {
	a := App{}
	a.Initialize("root", "12345678", "NursingHomes")

	a.Run(":8087")
}
