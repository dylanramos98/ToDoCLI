package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"time"
)

type Task struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	Status      bool      `json:"status"`
	Created     time.Time `json:"created"`
}

func main() {
	var task Task
	var tasks []Task

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
				Placeholder("Any extra information?").
				Value(&task.Description),

			// Prompt for tag
			huh.NewSelect[string]().
				Title("Choose a tag:").
				Options(
					huh.NewOptions("Personal", "School", "Work", "Other")...).
				Value(&task.Tag),
		),
	).
		WithTheme(huh.ThemeCharm())

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	task.Created = time.Now()

	// Add the new task to the task list
	tasks = append(tasks, task)

	// Save tasks to a JSON file
	filename := "tasks.json"
	err = SaveTasks(tasks, filename)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
	} else {
		fmt.Println("All set!")
	}
}
