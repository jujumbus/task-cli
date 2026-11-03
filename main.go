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
			listTasks, err := os.ReadFile("tasks.json")
			if err != nil {
				if !os.IsNotExist(err) {
					fmt.Fprintln(os.Stderr, "Ошибка при чтении файла: ", err)
					os.Exit(1)
				}
			} else {
				err = json.Unmarshal(listTasks, &allTasks)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Ошибка при разборе JSON", err)
					os.Exit(1)
				}
			}

			maxID := 0
			for _, task := range allTasks {
				if task.ID > maxID {
					maxID = task.ID
				}
			}

			description = strings.Join(os.Args[2:], " ")
			task := Task{ID: maxID + 1, Description: description, Status: "Выполнить", CreatedAt: now, UpdatedAt: now}
			allTasks = append(allTasks, task)

			dataTask, err := json.MarshalIndent(allTasks, "", "  ")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка кодировки задания", err)
				os.Exit(1)
			}

			err = os.WriteFile("tasks.json", dataTask, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Создание задания не удалось", err)
				os.Exit(1)
			}
			fmt.Println("Задание успешно сохранено")
		}

	case "list":
		{
			listTasks, err := os.ReadFile("tasks.json")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при чтении файла", err)
				os.Exit(1)
			}

			err = json.Unmarshal(listTasks, &allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при разборе задач списка")
				os.Exit(1)
			}

			if len(os.Args) >= 3 && os.Args[2] == "done" {
				fmt.Println("Выполненные задачи:\n")
				for _, task := range allTasks {
					if task.Status == "Выполнено" {
						fmt.Printf("Номер задачи: %d\nЗадача: %s\n\n", task.ID, task.Description)
					}
				}
			} else if len(os.Args) >= 3 && os.Args[2] == "progress" {
				fmt.Println("Задачи в процессе:\n")
				for _, task := range allTasks {
					if task.Status == "В процессе" {
						fmt.Printf("Номер задачи: %d\nЗадача: %s\n\n", task.ID, task.Description)
					}
				}
			} else if len(os.Args) >= 3 && os.Args[2] == "todo" {
				fmt.Println("Задачи к которым ещё не приступил:\n")
				for _, task := range allTasks {
					if task.Status == "Выполнить" {
						fmt.Printf("Номер задачи: %d\nЗадача: %s\n\n", task.ID, task.Description)
					}
				}
			} else {
				for _, task := range allTasks {
					fmt.Printf("Номер задачи: %d\nЗадача: %s\nСтатус: %s\n\n", task.ID, task.Description, task.Status)
				}
			}
		}

	case "update":
		{
			if len(os.Args) < 4 {
				fmt.Fprintln(os.Stderr, "Задание не может быть пустым")
				os.Exit(1)
			}
			listTasks, err := os.ReadFile("tasks.json")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при чтении файла: ", err)
				os.Exit(1)
			}
			err = json.Unmarshal(listTasks, &allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при разборе JSON", err)
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

			dataTask, err := json.MarshalIndent(allTasks, "", "  ")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка кодировки задания", err)
				os.Exit(1)
			}
			err = os.WriteFile("tasks.json", dataTask, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при изменении задачи", err)
				os.Exit(1)
			}
			fmt.Println("Задача обновлена успешно")
		}

	case "delete":
		{
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Не введён ID задачи")
				os.Exit(1)
			}

			listTasks, err := os.ReadFile("tasks.json")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при чтении файла", err)
				os.Exit(1)
			}

			err = json.Unmarshal(listTasks, &allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при разборе JSON", err)
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
			dataTask, err := json.MarshalIndent(allTasks, "", "  ")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка кодировки задания", err)
				os.Exit(1)
			}
			err = os.WriteFile("tasks.json", dataTask, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при изменении задачи", err)
				os.Exit(1)
			}
			fmt.Println("Задача удалена успешно")
		}

	case "done":
		{
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Не введён ID задачи")
				os.Exit(1)
			}

			listTasks, err := os.ReadFile("tasks.json")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при чтении файла: ", err)
				os.Exit(1)
			}
			err = json.Unmarshal(listTasks, &allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при разборе JSON", err)
				os.Exit(1)
			}

			taskID, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Неверно введён ID")
				os.Exit(1)
			}

			found := false
			for i, task := range allTasks {
				if taskID == task.ID {
					found = true
					taskPtr := &allTasks[i]
					taskPtr.Status = "Выполнено"
					taskPtr.UpdatedAt = now
				}
			}
			if !found {
				fmt.Fprintln(os.Stderr, "Задача с таким ID не найдена")
				os.Exit(1)
			}

			dataTask, err := json.MarshalIndent(allTasks, "", "  ")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка кодировки задания", err)
				os.Exit(1)
			}
			err = os.WriteFile("tasks.json", dataTask, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при изменении статуса задачи", err)
				os.Exit(1)
			}
			fmt.Println("Статус задачи обновлен успешно")
		}

	case "progress":
		{
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Не введён ID задачи")
				os.Exit(1)
			}

			listTasks, err := os.ReadFile("tasks.json")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при чтении файла", err)
				os.Exit(1)
			}
			err = json.Unmarshal(listTasks, &allTasks)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при разборе JSON", err)
				os.Exit(1)
			}

			taskID, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Неверно введён ID")
				os.Exit(1)
			}

			found := false
			for i, task := range allTasks {
				if taskID == task.ID {
					found = true
					taskPtr := &allTasks[i]
					taskPtr.Status = "В процессе"
					taskPtr.UpdatedAt = now
				}
			}
			if !found {
				fmt.Fprintln(os.Stderr, "Задача с таким ID не найдена")
				os.Exit(1)
			}

			dataTask, err := json.MarshalIndent(allTasks, "", "  ")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка кодирования задачи", err)
				os.Exit(1)
			}
			err = os.WriteFile("tasks.json", dataTask, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка при изменении статуса задачи", err)
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
