package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type Task struct {
	Id          int
	Title       string
	Description string
	Tag         string
	Status      bool
	Created     time.Time
}

func main() {
	fmt.Println("Let's build some tasks!")
	var tasks []Task
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Prompt user for input
		fmt.Println("\nEnter task details:")

		// Get Title
		fmt.Print("Title: ")
		scanner.Scan()
		title := scanner.Text()

		// Get Description
		fmt.Print("Description: ")
		scanner.Scan()
		description := scanner.Text()

		// Get Tag
		fmt.Print("Tag: ")
		scanner.Scan()
		tag := scanner.Text()

		// Add Task to List
		task := Task{
			Id:          len(tasks) + 1,
			Title:       title,
			Description: description,
			Tag:         tag,
			Status:      false,
			Created:     time.Now(),
		}
		tasks = append(tasks, task)

		// Ask if user wants to create another task
		fmt.Print("\nDo you want to create another task? (y/n): ")
		scanner.Scan()
		response := strings.ToLower(scanner.Text())
		if response != "y" {
			break
		}
	}

	// Display all tasks
	fmt.Println("\nAll Created Tasks:")
	for _, task := range tasks {
		fmt.Printf("\nTask #%d\nTitle: %s\nDescription: %s\nTag: %s\nCreated: %s\n",
			task.Id, task.Title, task.Description, task.Tag, humanize.Time(task.Created))
	}

}
