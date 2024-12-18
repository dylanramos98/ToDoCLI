package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"log"
	"os"
	"strings"
	"time"
)

type Task struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	Completed   bool      `json:"completed"`
	Created     time.Time `json:"created"`
}

const asciiArt = `
 /$$$$$$$$        /$$$$$$$             /$$$$$$  /$$       /$$$$$$
|__  $$__/       | $$__  $$           /$$__  $$| $$      |_  $$_/
   | $$  /$$$$$$ | $$  \ $$  /$$$$$$ | $$  \__/| $$        | $$  
   | $$ /$$__  $$| $$  | $$ /$$__  $$| $$      | $$        | $$  
   | $$| $$  \ $$| $$  | $$| $$  \ $$| $$      | $$        | $$  
   | $$| $$  | $$| $$  | $$| $$  | $$| $$    $$| $$        | $$  
   | $$|  $$$$$$/| $$$$$$$/|  $$$$$$/|  $$$$$$/| $$$$$$$$ /$$$$$$
   |__/ \______/ |_______/  \______/  \______/ |________/|______/
`

func main() {
	var tasks []Task
	var menuItem string
	for {
		// Create the main menu form
		menuForm := huh.NewForm(
			huh.NewGroup(huh.NewSelect[string]().
				Title(asciiArt).
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
		err := menuForm.Run()
		if err != nil {
			fmt.Println("Error displaying menu:", err)
			os.Exit(1)
		}

		// Handle the selected option
		switch menuItem {
		case "View Tasks":
			// Code to view tasks
			var filter string
			var title string

			viewForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Options(huh.NewOptions("Personal", "School", "Work", "Other")...).
						Title("Filter").
						Value(&filter),

					huh.NewSelect[string]().
						Value(&title).
						Height(8).
						TitleFunc(func() string {
							switch filter {
							default:
								return "Current Tasks"
							}
						}, &filter).
						OptionsFunc(func() []huh.Option[string] {
							opts := fetchTasksForFilter(filter, "tasks.json")
							return huh.NewOptions(opts...)
						}, &filter),
				),
			)
			err := viewForm.Run()
			if err != nil {
				log.Fatal(err)
			}

		case "Add Task":
			// Add a new task
			var task Task

			// Create a form for task creation
			taskForm := huh.NewForm(
				huh.NewGroup(
					// Prompt for task name
					huh.NewInput().
						Title(" Task name:").
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

func fetchTasksForFilter(filter string, filename string) []string {
	allTasks, err := ReadTasks(filename)
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return []string{}
	}

	// Collect task titles that match the filter
	var filteredTasks []string
	for _, task := range allTasks {
		if strings.EqualFold(task.Tag, filter) {
			filteredTasks = append(filteredTasks, task.Title)
		}
	}

	return filteredTasks
}
