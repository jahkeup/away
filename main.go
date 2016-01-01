package main

import (
	"fmt"
)

func main() {
	source := "."
	dest := "/tmp/dest"
	plan, _ := NewPlan(source, dest, PlanOptions{})
	plan.FindNodes()
	fmt.Println(Stowaway{}.Describe(plan))
}
