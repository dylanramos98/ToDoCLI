package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadTasks reads tasks from a JSON file. If the file doesn't exist, it returns an empty slice.
func ReadTasks(filename string) ([]Task, error) {
	// Initialize a slice for existing tasks
	var tasks []Task

	// Try to open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		// If the file doesn't exist, return an empty slice
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode tasks if the file opens successfully
	if err := json.NewDecoder(file).Decode(&tasks); err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// WriteTasks writes the given tasks to a JSON file, overwriting any existing content.
func WriteTasks(tasks []Task, filename string) error {
	// Reopen the file for writing (or create it if it doesn't exist)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write tasks to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	if err := encoder.Encode(tasks); err != nil {
		return fmt.Errorf("failed to write tasks: %w", err)
	}

	return nil
}

// SaveTasks reads existing tasks, appends new ones, and writes them back to the file.
func SaveTasks(newTasks []Task, filename string) error {
	// Read existing tasks from the file
	existingTasks, err := ReadTasks(filename)
	if err != nil {
		return err
	}

	// Append new tasks to the existing ones
	existingTasks = append(existingTasks, newTasks...)

	// Write the updated tasks back to the file
	return WriteTasks(existingTasks, filename)
}
