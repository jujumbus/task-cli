package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const (
	StatusDone     = "Выполнено"
	StatusProgress = "В процессе"
	StatusToDo     = "Выполнить"
)

func printTask(task Task) {
	fmt.Printf("Номер задачи: %d\nЗадача: %s\nСтатус: %s\nСоздана: %s\nОбновлена: %s\n\n", task.ID, task.Description, task.Status, task.CreatedAt.Format("02.01.2006 15:04"), task.UpdatedAt.Format("02.01.2006 15:04"))
}

func loadTasks() ([]Task, error) {
	listTasks, err := os.ReadFile("tasks.json")
	allTasks := []Task{}

	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}
	err = json.Unmarshal(listTasks, &allTasks)
	if err != nil {
		return nil, err
	}
	return allTasks, nil
}

func saveTasks(allTasks []Task) error {
	dataTask, err := json.MarshalIndent(allTasks, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("tasks.json", dataTask, 0644)
	if err != nil {
		return err
	}
	return nil
}

func changeStatus(allTasks []Task, taskID int, newStatus string) bool {
	for i, task := range allTasks {
		if taskID == task.ID {
			taskPtr := &allTasks[i]
			taskPtr.Status = newStatus
			taskPtr.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: task-cli <command>")
		os.Exit(1)
	}

	command := os.Args[1]
	description := ""
	now := time.Now()
	allTasks := []Task{}

	switch command {
	case "add":
		{
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Задание не может быть пустым")
				os.Exit(1)
			}
			var err error
			allTasks, err = loadTasks()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка загрузки задач:", err)
				os.Exit(1)
			}

			maxID := 0
			for _, task := range allTasks {
				if task.ID > maxID {
					maxID = task.ID
				}
			}

			description = strings.Join(os.Args[2:], " ")
			task := Task{ID: maxID + 1, Description: description, Status: StatusToDo, CreatedAt: now, UpdatedAt: now}
			allTasks = append(allTasks, task)

			err = saveTasks(allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка сохранения задачи:", err)
				os.Exit(1)
			}
			fmt.Println("Задание успешно сохранено")
		}

	case "list":
		{
			var err error
			allTasks, err = loadTasks()

			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка загрузки задач", err)
				os.Exit(1)
			}

			if len(os.Args) == 3 && os.Args[2] == "done" {
				fmt.Println("Выполненные задачи:\n")
				for _, task := range allTasks {
					if task.Status == StatusDone {
						printTask(task)
					}
				}
			} else if len(os.Args) == 3 && os.Args[2] == "progress" {
				fmt.Println("Задачи в процессе:\n")
				for _, task := range allTasks {
					if task.Status == StatusProgress {
						printTask(task)
					}
				}
			} else if len(os.Args) == 3 && os.Args[2] == "todo" {
				fmt.Println("Задачи к которым ещё не приступил:\n")
				for _, task := range allTasks {
					if task.Status == StatusToDo {
						printTask(task)
					}
				}
			} else if len(os.Args) == 2 {
				for _, task := range allTasks {
					printTask(task)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Неизвестный фильтр")
				os.Exit(1)
			}
		}

	case "update":
		{
			if len(os.Args) < 4 {
				fmt.Fprintln(os.Stderr, "Задание не может быть пустым")
				os.Exit(1)
			}
			var err error
			allTasks, err = loadTasks()

			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка загрузки задач", err)
				os.Exit(1)
			}

			taskID, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Неверный ID задачи")
				os.Exit(1)
			}

			found := false
			for i, task := range allTasks {
				if task.ID == taskID {
					found = true
					taskPtr := &allTasks[i]
					taskPtr.Description = strings.Join(os.Args[3:], " ")
					taskPtr.UpdatedAt = now
				}
			}
			if !found {
				fmt.Fprintln(os.Stderr, "Задача с таким ID не найдена")
				os.Exit(1)
			}

			err = saveTasks(allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка сохранения задачи:", err)
				os.Exit(1)
			}
			fmt.Println("Задача обновлена успешно")
		}

	case "delete":
		{
			if len(os.Args) != 3 {
				fmt.Fprintln(os.Stderr, "Проверьте формат команды: укажите ровно один ID и не добавляйте лишних аргументов")
				os.Exit(1)
			}

			var err error
			allTasks, err = loadTasks()

			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка загрузки задач", err)
				os.Exit(1)
			}

			taskID, err := strconv.Atoi(os.Args[2])
			remainingTasks := []Task{}
			found := false
			if err != nil {
				fmt.Fprintln(os.Stderr, "Неверно введён ID")
				os.Exit(1)
			}
			for i, task := range allTasks {
				if taskID == task.ID {
					found = true
				} else {
					remainingTasks = append(remainingTasks, allTasks[i])
				}
			}

			if !found {
				fmt.Fprintln(os.Stderr, "Задача с таким ID не найдена")
				os.Exit(1)
			}

			allTasks = remainingTasks
			err = saveTasks(allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка сохранения задачи:", err)
				os.Exit(1)
			}
			fmt.Println("Задача удалена успешно")
		}

	case "done":
		{
			if len(os.Args) != 3 {
				fmt.Fprintln(os.Stderr, "Проверьте формат команды: укажите ровно один ID и не добавляйте лишних аргументов")
				os.Exit(1)
			}

			var err error
			allTasks, err = loadTasks()

			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка загрузки задач", err)
				os.Exit(1)
			}

			taskID, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Неверно введён ID")
				os.Exit(1)
			}

			found := changeStatus(allTasks, taskID, StatusDone)
			if !found {
				fmt.Fprintln(os.Stderr, "Задача с таким ID не найдена")
				os.Exit(1)
			}

			err = saveTasks(allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка сохранения задачи:", err)
				os.Exit(1)
			}
			fmt.Println("Статус задачи обновлен успешно")
		}

	case "progress":
		{
			if len(os.Args) != 3 {
				fmt.Fprintln(os.Stderr, "Проверьте формат команды: укажите ровно один ID и не добавляйте лишних аргументов")
				os.Exit(1)
			}

			var err error
			allTasks, err = loadTasks()

			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка загрузки задач", err)
				os.Exit(1)
			}

			taskID, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Неверно введён ID")
				os.Exit(1)
			}

			found := changeStatus(allTasks, taskID, StatusProgress)
			if !found {
				fmt.Fprintln(os.Stderr, "Задача с таким ID не найдена")
				os.Exit(1)
			}

			err = saveTasks(allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка сохранения задачи:", err)
				os.Exit(1)
			}
			fmt.Println("Статус задачи обновлён успешно")
		}

	default:
		{
			fmt.Fprintln(os.Stderr, "Неизвестная команда:", command)
			os.Exit(1)
		}

	}
}
