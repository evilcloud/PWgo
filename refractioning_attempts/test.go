package main

import "fmt"

type Employee struct {
	Name string
	Age  int
}

func (obj *Employee) Info() {
	if obj.Name == "" {
		obj.Name = "John Doe"
	}
	if obj.Age == 0 {
		obj.Age = 25
	}
}

func main() {
	emp1 := Employee{Name: "Mr. Fred"}
	emp1.Info()
	fmt.Println("1:", emp1)

	emp2 := Employee{Age: 26}
	emp2.Info()
	fmt.Println("2:", emp2)

	emp3 := Employee{}
	emp3.Name = "Mrs. Smith"
	emp3.Age = 20
	fmt.Println("3:", emp3)

	emp4 := Employee{}
	emp4.Info()
	fmt.Println("4:", emp4)
}
