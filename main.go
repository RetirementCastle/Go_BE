// main.go

package main

func main() {
	a := App{}
	a.Initialize("root", "mi nombre es 123", "NursingHomes")

	a.Run(":8080")
}
