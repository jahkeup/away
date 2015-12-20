package main

type AwayOptions struct{}

type Away struct {
	Options AwayOptions

	Target string
	Source string

	Plan *Plan
}

type Plan struct {
}
