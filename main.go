// main.go

package main

func main() {
	a := App{}
	//a.Initialize("root", "12345678", "NursingHomes")
	a.Initialize("generic_test", "Arquitectura2018", "retirementcastle")
	a.Run(":8087")
}
