package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	todo "github.com/1set/todotxt"
	"github.com/TwinProduction/go-color"
	"github.com/jedib0t/go-pretty/table"
)

func main() {

	// Create file if not exist

	if _, err := os.Stat("todo.txt"); os.IsNotExist(err) {
		os.Create("todo.txt")
	}
	// Load and perform operations on the tasklist
	if todolist, err := todo.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	} else {
		inputArgs := os.Args[1:]
		resultRegex := formatArgs(inputArgs)

		if inputArgs[0] == "ls" || inputArgs[0] == "completed" {
			handlelscompleted(todolist, inputArgs, resultRegex)
		} else if inputArgs[0] == "add" {
			handleAdd(todolist, inputArgs)
		} else if inputArgs[0] == "rm" {
			taskID, _ := strconv.Atoi(inputArgs[1])
			handleRm(todolist, taskID)
		} else if inputArgs[0] == "do" {
			taskID, _ := strconv.Atoi(inputArgs[1])
			handleDo(todolist, taskID)
		} else if inputArgs[0] == "tags" {
			handleTags(todolist)
		} else if inputArgs[0] == "projects" {
			handleProjects(todolist)
		} else if inputArgs[0] == "due" {
			handleDue(todolist, inputArgs)
		} else if inputArgs[0] == "extend" {
			handleExtend(todolist, inputArgs)
		} else if inputArgs[0] == "help" {
			handleHelp(inputArgs)
		}

	}
}

/*
 *	Keep a track of regexes of the parameters
 */

func formatArgs(params []string) []string {
	// fmt.Println(params)
	var resultRegex []string

	for _, param := range params {
		if val, _ := regexp.MatchString("^[|]?@[a-zA-Z0-9]+$", param); val == true {
			// context
			resultRegex = append(resultRegex, "^[|]?@[a-zA-Z0-9]+$")
		} else if val, _ := regexp.MatchString("^[|]?\\+[a-zA-Z0-9]+$", param); val == true {
			// project
			resultRegex = append(resultRegex, "^[|]?\\+[a-zA-Z0-9]+$")
		} else if val, _ := regexp.MatchString("^[|]?\\([a-zA-Z0-9]+\\)+$", param); val == true {
			// priority
			resultRegex = append(resultRegex, "^[|]?\\([a-zA-Z0-9]+\\)+$")
		} else if val, _ := regexp.MatchString("^[|]?[a-zA-Z0-9]+:$", param); val == true {
			// tag
			resultRegex = append(resultRegex, "^[|]?[a-zA-Z0-9]+:$")
		} else if val, _ := regexp.MatchString("^[|]?>\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$", param); val == true {
			// after date
			resultRegex = append(resultRegex, "^[|]?>\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$")
		} else if val, _ := regexp.MatchString("^[|]?<\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$", param); val == true {
			// before date
			resultRegex = append(resultRegex, "^[|]?<\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$")
		} else if val, _ := regexp.MatchString("[|]?completed", param); val == true {
			resultRegex = append(resultRegex, "[|]?completed")
		}

	}

	return resultRegex

}

/*
 *	Functions to handle various parameters
 */

func handlelscompleted(todolist todo.TaskList, inputArgs []string, resultRegex []string) {
	// If no other paramters present, print all.
	var hasCompleted bool
	var onlyCompleted bool
	var hasSort string
	if len(inputArgs) == 1 {
		if inputArgs[0] == "completed" {
			todolist = todolist.Filter(todo.FilterCompleted)
			todolist.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc)
			onlyCompleted = true
		} else {
			todolist = todolist.Filter(todo.FilterNotCompleted)
			todolist.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc)
		}
	} else {
		// Check if matches any context
		var isContextChanged bool = false
		var contextList todo.TaskList
		for _, param := range inputArgs {
			if val, _ := regexp.MatchString("^@[a-zA-Z0-9]+$", param); val == true {
				isContextChanged = true
				param = param[1:]
				pred := todo.FilterByContext(param)
				newList := todolist.Filter(pred)
				for _, newTask := range newList {
					var isExist bool = false
					for _, oldTask := range contextList {
						if newTask.Original == oldTask.Original {
							isExist = true
						}
					}
					if isExist == false {
						contextList.AddTask(&newTask)
					}
				}
			}
		}
		if isContextChanged {
			todolist = contextList
		}

		// Check if matches any projects
		var projectList todo.TaskList
		var isProjectChanged bool = false
		for _, param := range inputArgs {
			if val, _ := regexp.MatchString("^\\+[a-zA-Z0-9]+$", param); val == true {
				isProjectChanged = true
				param = param[1:]
				pred := todo.FilterByProject(param)
				newList := todolist.Filter(pred)
				for _, newTask := range newList {
					var isExist bool = false
					for _, oldTask := range projectList {
						if newTask.Original == oldTask.Original {
							isExist = true
						}
					}
					if isExist == false {
						projectList.AddTask(&newTask)
					}
				}

			}
		}
		if isProjectChanged {
			todolist = projectList
		}

		// Check if any priorities match
		var priorityList todo.TaskList
		var isPriorityChanged bool = false
		for _, param := range inputArgs {
			if val, _ := regexp.MatchString("^\\([a-zA-Z0-9]+\\)+$", param); val == true {

				isPriorityChanged = true
				param := param[1 : len(param)-1]
				pred := todo.FilterByPriority(param)
				newList := todolist.Filter(pred)
				for _, newTask := range newList {
					var isExist bool = false
					for _, oldTask := range priorityList {
						if newTask.Original == oldTask.Original {
							isExist = true
						}
					}
					if isExist == false {
						priorityList.AddTask(&newTask)
					}
				}

			}
		}
		if isPriorityChanged {
			todolist = priorityList
		}

		// Check if any tags match
		var tagList todo.TaskList
		var isTagsListChanged bool = false
		for _, param := range inputArgs {
			if val, _ := regexp.MatchString("^[a-zA-Z0-9]+:$", param); val == true {

				isTagsListChanged = true
				param := param[:len(param)-1]
				newList := FilterByAdditionalTags(todolist, param)

				for _, newTask := range newList {
					var isExist bool = false
					for _, oldTask := range tagList {
						if newTask.Original == oldTask.Original {
							isExist = true
						}
					}
					if isExist == false {
						tagList.AddTask(&newTask)
					}
				}

			}
		}
		if isTagsListChanged {
			todolist = tagList
		}
		// Check for specified date constraints, "completed" param, and order param
		for _, param := range inputArgs {
			if val, _ := regexp.MatchString("^>\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$", param); val == true {
				// after date
				resultRegex = append(resultRegex, "^>\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$")
				param := param[1:]
				newTime, _ := time.Parse("2006-01-02", param)
				for l := 0; l < len(todolist); l++ {
					task := todolist[l]
					var dateVar time.Time
					if onlyCompleted {
						dateVar = task.DueDate
					} else {
						for _, seg := range task.Segments() {
							if seg.Type == todo.SegmentCompletedDate {
								dateVar, _ = time.Parse("2006-01-02", seg.Display)
							}
						}
					}
					if dateVar.Before(newTime) {
						// Remove task from todolist
						todolist = append(todolist[:l], todolist[l+1:]...)
						l--
					} else {
					}
				}
			} else if val, _ := regexp.MatchString("^<\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$", param); val == true {
				// before date
				resultRegex = append(resultRegex, "^<\\d{4}\\-(0[1-9]|1[012])\\-(0[1-9]|[12][0-9]|3[01])+$")
				param := param[1:]
				newTime, _ := time.Parse("2006-01-02", param)
				for l := 0; l < len(todolist); l++ {
					task := todolist[l]
					var dateVar time.Time
					if onlyCompleted {
						dateVar = task.DueDate
					} else {
						for _, seg := range task.Segments() {
							if seg.Type == todo.SegmentCompletedDate {
								dateVar, _ = time.Parse("2006-01-02", seg.Display)
							}
						}
					}
					if dateVar.After(newTime) {
						// Remove task from todolist
						todolist = append(todolist[:l], todolist[l+1:]...)
						l--
					}
				}
			} else if val, _ := regexp.MatchString("completed", param); val == true {
				resultRegex = append(resultRegex, "completed")
				hasCompleted = true
			} else if val, _ := regexp.MatchString("\\|[a-zA-Z0-9]+", param); val == true {
				resultRegex = append(resultRegex, "\\|[a-zA-Z0-9]+")
				hasSort = param[1:]
			}
		}
	}
	var printComplete bool
	var completedTasks todo.TaskList
	var incomplete todo.TaskList
	if hasCompleted && inputArgs[0] == "ls" {
		// if it contains "completed"
		printComplete = true
		todolist.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc)
		incomplete = todolist.Filter(todo.FilterNotCompleted)
		completedTasks = todolist.Filter(todo.FilterCompleted)
	} else {
		if inputArgs[0] == "completed" {
			todolist = todolist.Filter(todo.FilterCompleted)
			todolist.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc)
		} else {
			todolist = todolist.Filter(todo.FilterNotCompleted)
			todolist.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc)
		}
	}

	if len(hasSort) != 0 {
		if hasSort == "TaskIDAsc" {
			todolist.Sort(todo.SortTaskIDAsc)
		} else if hasSort == "TaskIDDesc" {
			todolist.Sort(todo.SortTaskIDDesc)
		} else if hasSort == "TodoTextAsc" {
			todolist.Sort(todo.SortTodoTextAsc)
		} else if hasSort == "TodoTextDesc" {
			todolist.Sort(todo.SortTodoTextDesc)
		} else if hasSort == "PriorityAsc" {
			todolist.Sort(todo.SortPriorityAsc)
		} else if hasSort == "PriorityDesc" {
			todolist.Sort(todo.SortPriorityDesc)
		} else if hasSort == "CreatedDateAsc" {
			todolist.Sort(todo.SortCreatedDateAsc)
		} else if hasSort == "CreatedDateDesc" {
			todolist.Sort(todo.SortCreatedDateDesc)
		} else if hasSort == "CompletedDateAsc" {
			todolist.Sort(todo.SortCompletedDateAsc)
		} else if hasSort == "CompletedDateDesc" {
			todolist.Sort(todo.SortCompletedDateDesc)
		} else if hasSort == "DueDateAsc" {
			todolist.Sort(todo.SortDueDateAsc)
		} else if hasSort == "DueDateDesc" {
			todolist.Sort(todo.SortDueDateDesc)
		} else if hasSort == "ContextAsc" {
			todolist.Sort(todo.SortContextAsc)
		} else if hasSort == "ContextDesc" {
			todolist.Sort(todo.SortContextDesc)
		} else if hasSort == "ProjectAsc" {
			todolist.Sort(todo.SortPriorityAsc)
		} else if hasSort == "ProjectDesc" {
			todolist.Sort(todo.SortProjectDesc)
		}
	}
	if printComplete {

		for _, task := range incomplete {
			if task.Priority == "A" {
				fmt.Println(color.Ize(color.Yellow, task.Original))
			} else if task.Priority == "B" {
				fmt.Println(color.Ize(color.Red, task.Original))
			} else if task.Priority == "C" {
				fmt.Println(color.Ize(color.Green, task.Original))
			} else if task.Priority == "D" {
				fmt.Println(color.Ize(color.Cyan, task.Original))
			} else if task.Priority == "E" {
				fmt.Println(color.Ize(color.Blue, task.Original))
			} else {
				fmt.Println(task.Original)
			}
		}
		for _, task := range completedTasks {
			if task.Priority == "A" {
				fmt.Println(color.Ize(color.Yellow, task.Original))
			} else if task.Priority == "B" {
				fmt.Println(color.Ize(color.Red, task.Original))
			} else if task.Priority == "C" {
				fmt.Println(color.Ize(color.Green, task.Original))
			} else if task.Priority == "D" {
				fmt.Println(color.Ize(color.Cyan, task.Original))
			} else if task.Priority == "E" {
				fmt.Println(color.Ize(color.Blue, task.Original))
			} else {
				fmt.Println(task.Original)
			}
		}
		// fmt.Println(incomplete)
		// fmt.Print(completedTasks)
		fmt.Printf("\nTOTAL:{%d}\n", len(incomplete)+len(completedTasks))
	} else {
		// fmt.Println(todolist)
		for _, task := range todolist {
			if task.Priority == "A" {
				fmt.Println(color.Ize(color.Yellow, task.Original))
			} else if task.Priority == "B" {
				fmt.Println(color.Ize(color.Red, task.Original))
			} else if task.Priority == "C" {
				fmt.Println(color.Ize(color.Green, task.Original))
			} else if task.Priority == "D" {
				fmt.Println(color.Ize(color.Cyan, task.Original))
			} else if task.Priority == "E" {
				fmt.Println(color.Ize(color.Blue, task.Original))
			} else {
				fmt.Println(task.Original)
			}
		}
		fmt.Printf("TOTAL:{%d}\n", len(todolist))
	}

}

// FilterByAdditionalTags - Custom function to filter by additional tags that returns a new todo.TaskList object
func FilterByAdditionalTags(todolist todo.TaskList, newTag string) todo.TaskList {
	var newTaskList todo.TaskList

	for _, task := range todolist {
		if _, ok := task.AdditionalTags[newTag]; ok {
			// If not in newTaskList, add it
			var isExist bool = false
			for _, newTask := range newTaskList {
				if newTask.Original == task.Original {
					isExist = true
				}
			}
			if isExist == false {
				newTaskList.AddTask(&task)
			}
		}
	}

	return newTaskList
}

// handleAdd - Adds a new task
func handleAdd(todolist todo.TaskList, inputArgs []string) {
	taskString := inputArgs[1]
	newTask, _ := todo.ParseTask(taskString)
	todolist = append(todolist, *newTask)
	if err := todo.WriteToPath(&todolist, "todo.txt"); err != nil {
		log.Fatal(err)
	}
}

// handleRm - Remove a task
func handleRm(todolist todo.TaskList, taskID int) {
	if err := todolist.RemoveTaskByID(taskID); err != nil {
		log.Fatal(err)
	}
	if err := todo.WriteToPath(&todolist, "todo.txt"); err != nil {
		log.Fatal(err)
	}
}

// handleDo - Mark a task as completed
func handleDo(todolist todo.TaskList, taskID int) {
	for i, task := range todolist {
		if task.ID == taskID {
			todolist[i].Complete()
		}
	}
	if err := todo.WriteToPath(&todolist, "todo.txt"); err != nil {
		log.Fatal(err)
	}
}

// handleTags - List all the tags in the tasks (no duplicates)
func handleTags(todolist todo.TaskList) {
	// Use the keys of a map as the storing variable to avoid duplicates
	tagMap := make(map[string]int, 0)

	for _, task := range todolist {
		for tag, _ := range task.AdditionalTags {
			tagMap[tag] = 1
		}
	}

	for tag, _ := range tagMap {
		fmt.Println(tag)
	}
}

func handleProjects(todolist todo.TaskList) {
	// Use the keys of a map as the storing variable to avoid duplicates
	projectMap := make(map[string]int, 0)

	for _, task := range todolist {

		for _, project := range task.Projects {
			projectMap[project] = 1
		}
	}

	for project, _ := range projectMap {
		fmt.Println(project)
	}
}

func handleDue(todolist todo.TaskList, inputArgs []string) {
	dueMap := make(map[string][]int, 0)
	noDueDate := []int{0, 0}
	var dueLessThan bool = false
	var dueMoreThan bool = false
	var dateConstraintLessThan time.Time
	var dateConstraintMoreThan time.Time
	if len(inputArgs) != 1 {

		for _, param := range inputArgs {
			if strings.HasPrefix(param, "<=") {
				dueLessThan = true
				dateConstraintLessThan, _ = time.Parse("2006-01-02", param[2:])
			} else if strings.HasPrefix(param, ">=") {
				dueMoreThan = true
				dateConstraintMoreThan, _ = time.Parse("2006-01-02", param[2:])
			}
		}

	}
	for _, task := range todolist {
		// defaultDate, _ := time.Parse("2006-01-02", "0001-01-01")

		if task.HasDueDate() == false {
			if !task.IsCompleted() {
				noDueDate[1]++
			}
			noDueDate[0]++
		} else {
			if dueLessThan && dueMoreThan {
				if task.DueDate.After(dateConstraintMoreThan) && task.DueDate.Before(dateConstraintLessThan) {

					timeFormatNew := task.DueDate.String()
					timeFormatNew = timeFormatNew[:len(timeFormatNew)-19]
					count := dueMap[timeFormatNew]
					if len(count) == 0 {
						count = []int{0, 0}
					}
					count[0] = count[0] + 1
					if !task.IsCompleted() {
						count[1] = count[1] + 1
					}
					dueMap[timeFormatNew] = count
				}
			} else if dueLessThan {
				if task.DueDate.Before(dateConstraintLessThan) {
					timeFormatNew := task.DueDate.String()
					timeFormatNew = timeFormatNew[:len(timeFormatNew)-19]
					count := dueMap[timeFormatNew]
					if len(count) == 0 {
						count = []int{0, 0}
					}
					count[0] = count[0] + 1
					if !task.IsCompleted() {
						count[1] = count[1] + 1
					}
					dueMap[timeFormatNew] = count
				}
			} else if dueMoreThan {
				if task.DueDate.After(dateConstraintMoreThan) {
					timeFormatNew := task.DueDate.String()
					timeFormatNew = timeFormatNew[:len(timeFormatNew)-19]
					count := dueMap[timeFormatNew]
					if len(count) == 0 {
						count = []int{0, 0}
					}
					count[0] = count[0] + 1
					if !task.IsCompleted() {
						count[1] = count[1] + 1
					}
					dueMap[timeFormatNew] = count
				}
			} else {
				timeFormatNew := task.DueDate.String()
				timeFormatNew = timeFormatNew[:len(timeFormatNew)-19]
				count := dueMap[timeFormatNew]
				if len(count) == 0 {
					count = []int{0, 0}
				}
				count[0] = count[0] + 1
				if !task.IsCompleted() {
					count[1] = count[1] + 1
				}
				dueMap[timeFormatNew] = count
			}
		}

	}

	for dueDate, dueCount := range dueMap {

		fmt.Printf("%s\t%d (%d)\n", dueDate, dueCount[0], dueCount[1])
	}
	fmt.Printf("%d (%d)\n", noDueDate[0], noDueDate[1])
}

func handleExtend(todolist todo.TaskList, inputArgs []string) {

	taskID, _ := strconv.Atoi(inputArgs[1])
	quantity, _ := strconv.Atoi(inputArgs[2])
	unit := inputArgs[3]

	task, _ := todolist.GetTask(taskID)

	if task.HasDueDate() == false {
		task.DueDate = time.Now()
	}

	if unit == "day" {
		task.DueDate = task.DueDate.AddDate(0, 0, quantity)
	} else if unit == "week" {
		task.DueDate = task.DueDate.AddDate(0, 0, 7*quantity)
	} else if unit == "month" {
		task.DueDate = task.DueDate.AddDate(0, quantity, 0)
	} else if unit == "year" {
		task.DueDate = task.DueDate.AddDate(quantity, 0, 0)
	}

	if err := todo.WriteToPath(&todolist, "todo.txt"); err != nil {
		log.Fatal(err)
	}
}

func handleHelp(inputArgs []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Command", "Description"})
	t.AppendRows([]table.Row{
		{"ls <options>", "Display the entire todo list. Options include @context, +project, (priority), tag:, <>datestring, |order, completed"},
		{"completed <(><)datestring>", "Show only the completed tasks"},
		{"add <task>", "Add a new task"},
		{"rm <taskID>", "Remove a task with the given ID"},
		{"do <taskID>", "Mark the task as completed"},
		{"tags", "Display all the unique tags"},
		{"projects", "Display all the unique projects"},
		{"due <options>", "Display all the unique due dates along with the number of projects due and the projects that remain incomplete to that day."},
		{"extend <taskID> <quantity> <unit>", "Extend the due date of a task by the given quantity. The units available are day, week, month, year"},
	})

	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Render()
}
