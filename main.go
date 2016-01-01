package main

import (
	"fmt"
)

func main() {
	source := "."
	dest := "/tmp/dest"
	plan, _ := NewPlan(source, dest, PlanOptions{
		LinkFilesOnly: true,
	})
	plan.FindNodes()
	fmt.Println(string(Stowaway{}.Describe(plan)))
}
