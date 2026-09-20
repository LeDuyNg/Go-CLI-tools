package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"todo"
)

// Default file name
var todoFileName = ".todo.json"

// Get vault path
func getVaultPath() (string, error) {
	if vault := os.Getenv("TODO_VAULT"); vault != "" {
		return vault, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".todo-vault"), nil
}

// getTask function decides where to get the description for a new
// task from: arguments or STDIN
func getTask(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}

	scanner := bufio.NewScanner(r)
	var tasks []string
	for scanner.Scan() {
		task := strings.TrimSpace(scanner.Text())

		if task == "" {
			continue // ignore blank lines
		}

		tasks = append(tasks, task)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("No tasks provided")
	}

	return tasks, nil
}

func main() {
	vaultDir, err := getVaultPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := os.MkdirAll(vaultDir, 0700); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	todoFilePath := filepath.Join(vaultDir, todoFileName)

	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the Todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("del", 0, "Delete an item from the list")
	verbose := flag.Bool("v", false, "Verbose Output")
	undone := flag.Bool("undone", false, "Flag to use with list, only show unfinished tasks")
	done := flag.Bool("done", false, "Flag to use with list, only show finished tasks")

	flag.Parse()
	l := &todo.List{}

	// Use the Get method to read to do items from file
	if err := l.Get(todoFilePath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	// For no extra arguments, print the list
	case *list && !*verbose && !*undone && !*done:
		// List current to do items
		fmt.Print(l)
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFilePath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// When any arguments (excluding flags) are provided, they will be
		// used as the new task
		tasks, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		for _, t := range tasks {
			l.Add(t)
		}

		// Save the new list
		if err := l.Save(todoFilePath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		// Delete the given item
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFilePath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *list && *verbose && !*undone && !*done: // list all, verbose
		fmt.Print(l.StringVerbose())
	case *list && !*verbose && *undone && !*done: // list all undone
		// Temporary list to contain undone tasks
		undoneList := &todo.List{}

		for _, t := range *l {
			if !t.Done {
				*undoneList = append(*undoneList, t)
			}
		}
		fmt.Print(undoneList)
	case *list && *verbose && *undone && !*done: // list all undone, verbose
		// Temporary list to contain undone tasks
		undoneList := &todo.List{}

		for _, t := range *l {
			if !t.Done {
				*undoneList = append(*undoneList, t)
			}
		}
		fmt.Print(undoneList.StringVerbose())
	case *list && !*verbose && !*undone && *done: // list all done
		// Temporary list to contain done tasks
		doneList := &todo.List{}

		for _, t := range *l {
			if t.Done {
				*doneList = append(*doneList, t)
			}
		}
		fmt.Print(doneList)
	case *list && *verbose && !*undone && *done: // list all done, verbose
		// Temporary list to contain undone tasks
		doneList := &todo.List{}

		for _, t := range *l {
			if t.Done {
				*doneList = append(*doneList, t)
			}
		}
		fmt.Print(doneList.StringVerbose())
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}
