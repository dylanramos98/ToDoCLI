package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"time"
)

type Task struct {
	Title       string
	Description string
	Tag         string
	Status      bool
	Created     time.Time
}

func main() {
	var task Task
	//var taskList []Task

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("ToDoCLI").
			Description("Welcome to _ToDoCLI!_.\n\n").
			Next(true).
			NextLabel("Proceed"),
		),

		// Create a Task
		huh.NewGroup(huh.NewInput().
			Title("Task name:").
			Prompt("* ").
			Validate(func(t string) error {
				if t == "" {
					return fmt.Errorf("please enter a task name")
				}
				return nil
			}).
			Value(&task.Title),
		),
	)
	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
