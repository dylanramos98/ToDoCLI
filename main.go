package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
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

		// Select a task to display
		switch menuItem {
		case "View Tasks":
			// Code to view tasks
			var filter string
			var title string
			var filteredTasks []Task

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
							titles, tasks, err := fetchTasksForFilter(filter, "tasks.json")
							if err != nil {
								log.Fatalf("Error fetching tasks: %v", err)
							}
							filteredTasks = tasks
							return huh.NewOptions(titles...)
						}, &filter),
				),
			)
			err := viewForm.Run()
			if err != nil {
				log.Fatal(err)
			}

			// Find and display the selected task
			var selectedTask *Task
			for _, task := range filteredTasks {
				if task.Title == title {
					selectedTask = &task
					break
				}
			}

			{
				var sb strings.Builder
				keyword := func(s string) string {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
				}
				fmt.Fprintf(&sb,
					"\n %s \n\n Title: \n %s \n Description: \n %s \n Tag: \n %s \n Created:\n %v \n",
					lipgloss.NewStyle().Bold(true).Render("Task"),
					keyword(selectedTask.Title),
					keyword(selectedTask.Description),
					keyword(selectedTask.Tag),
					keyword(humanize.Time(selectedTask.Created)),
				)

				fmt.Println(
					lipgloss.NewStyle().
						Width(40).
						BorderStyle(lipgloss.RoundedBorder()).
						BorderForeground(lipgloss.Color("63")).
						Padding(1, 2).
						Render(sb.String()),
				)
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

func fetchTasksForFilter(filter string, filename string) ([]string, []Task, error) {
	allTasks, err := ReadTasks(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading tasks: %w", err)
	}

	var filteredTasks []Task
	var filteredTitles []string
	for _, task := range allTasks {
		if strings.EqualFold(task.Tag, filter) {
			filteredTasks = append(filteredTasks, task)
			filteredTitles = append(filteredTitles, task.Title)
		}
	}

	return filteredTitles, filteredTasks, nil
}
