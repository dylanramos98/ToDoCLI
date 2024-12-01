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
			Title("\nToDoCLI").
			Description("Welcome to _ToDoCLI!_\n\n").
			Next(true).
			NextLabel("Proceed"),
		),

		// Create a Task
		huh.NewGroup(
			// Prompt for task name
			huh.NewInput().
				Title("Task name:").
				Prompt("* ").
				Validate(func(t string) error {
					if t == "" {
						return fmt.Errorf("please enter a task name")
					}
					return nil
				}).
				Value(&task.Title),

			// Prompt for description
			huh.NewText().
				Title("Description:").
				Value(&task.Description),

			huh.NewSelect[string]().
				Title("Choose a tag:").
				Options(
					huh.NewOptions("Personal", "School", "Work", "Other")...).
				Value(&task.Tag),
		),
	)
	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
