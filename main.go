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
	var tasks []Task
	var menuItem string
	for {
		// Create the main menu form
		form := huh.NewForm(
			huh.NewGroup(huh.NewSelect[string]().
				Title("\nWelcome to ToDoCLI!\n").
				Description("What would you like to do?\n").
				Options(
					huh.NewOption("View Tasks", "View Tasks"),
					huh.NewOption("Add Task", "Add Task"),
					huh.NewOption("Exit", "Exit"),
				).
				Value(&menuItem),
			),
		)

		// Run the form
		err := form.Run()
		if err != nil {
			fmt.Println("Error displaying menu:", err)
			os.Exit(1)
		}

		// Handle the selected option
		switch menuItem {
		case "View Tasks":
			// Code to view tasks

		case "Add Task":
			// Add a new task
			var task Task

			// Create a form for task creation
			taskForm := huh.NewForm(
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
			)

			err := taskForm.Run()
			if err != nil {
				fmt.Println("Error creating task:", err)
				continue
			}

			// Set the task creation time
			task.Created = time.Now()

			// Add the new task to the task list
			tasks = append(tasks, task)

			// Save tasks to a JSON file
			filename := "tasks.json"
			err = SaveTasks(tasks, filename)
			if err != nil {
				fmt.Println("Error saving task:", err)
			} else {
				fmt.Println("All set!")
			}
			tasks = []Task{}

		case "Exit":
			// Exit the loop
			fmt.Println("See ya later!")
			return

		default:
			fmt.Println("Invalid option selected. Please try again.")
		}
	}
}
