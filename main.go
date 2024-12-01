package main

import (
	"encoding/json"
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

// SaveTasks writes the given tasks to a JSON file.
// If the file exists, it reads and appends to the existing tasks.
func SaveTasks(newTasks []Task, filename string) error {
	// Initialize a slice for existing tasks
	var existingTasks []Task

	// Try to open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		// If file doesn't exist, assume no existing tasks
		if os.IsNotExist(err) {
			existingTasks = []Task{}
		} else {
			return fmt.Errorf("failed to open file: %w", err)
		}
	} else {
		// Decode existing tasks if the file opens successfully
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&existingTasks); err != nil && err.Error() != "EOF" {
			return fmt.Errorf("failed to decode existing tasks: %w", err)
		}
	}

	// Append new tasks to the existing ones
	existingTasks = append(existingTasks, newTasks...)

	// Reopen the file for writing (or create it if it doesn't exist)
	file, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the updated tasks back to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	if err := encoder.Encode(existingTasks); err != nil {
		return fmt.Errorf("failed to write tasks: %w", err)
	}

	return nil
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
		fmt.Println("Tasks saved to", filename)
	}
}
