package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

var scanner = bufio.NewScanner(os.Stdin)

// askString prompts for a string input
func askString(prompt string, defaultValue string) string {
    if defaultValue != "" {
        fmt.Printf("%s [%s]: ", prompt, defaultValue)
    } else {
        fmt.Printf("%s: ", prompt)
    }

    scanner.Scan()
    input := strings.TrimSpace(scanner.Text())

    if input == "" && defaultValue != "" {
        return defaultValue
    }

    return input
}

// askInt prompts for an integer input
func askInt(prompt string, defaultValue int) int {
    for {
        if defaultValue > 0 {
            fmt.Printf("%s [%d]: ", prompt, defaultValue)
        } else {
            fmt.Printf("%s: ", prompt)
        }

        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())

        if input == "" && defaultValue > 0 {
            return defaultValue
        }

        num, err := strconv.Atoi(input)
        if err != nil {
            fmt.Println("Please enter a valid number")
            continue
        }

        return num
    }
}

// askChoice prompts for a choice from options
func askChoice(prompt string, options []string, defaultIndex int) int {
    fmt.Println(prompt)
    for i, opt := range options {
        if i == defaultIndex {
            fmt.Printf("  %d) %s (default)\n", i+1, opt)
        } else {
            fmt.Printf("  %d) %s\n", i+1, opt)
        }
    }

    for {
        if defaultIndex >= 0 {
            fmt.Printf("Enter choice [%d]: ", defaultIndex+1)
        } else {
            fmt.Print("Enter choice: ")
        }

        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())

        if input == "" && defaultIndex >= 0 {
            return defaultIndex
        }

        choice, err := strconv.Atoi(input)
        if err != nil || choice < 1 || choice > len(options) {
            fmt.Printf("Please enter a number between 1 and %d\n", len(options))
            continue
        }

        return choice - 1
    }
}

// askMultiChoice prompts for multiple choices from options
func askMultiChoice(prompt string, options []string) []int {
    fmt.Println(prompt)
    for i, opt := range options {
        fmt.Printf("  %d) %s\n", i+1, opt)
    }
    fmt.Println("Enter choices separated by commas (e.g., 1,3,5) or 'done' when finished:")

    for {
        fmt.Print("Choices: ")
        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())

        if input == "" || strings.ToLower(input) == "done" {
            return []int{}
        }

        parts := strings.Split(input, ",")
        var choices []int
        valid := true

        for _, part := range parts {
            choice, err := strconv.Atoi(strings.TrimSpace(part))
            if err != nil || choice < 1 || choice > len(options) {
                fmt.Printf("Invalid choice: %s. Please enter numbers between 1 and %d\n", part, len(options))
                valid = false
                break
            }
            choices = append(choices, choice-1)
        }

        if valid {
            return choices
        }
    }
}

// askYesNo prompts for a yes/no answer
func askYesNo(prompt string, defaultYes bool) bool {
    suffix := " [Y/n]: "
    if !defaultYes {
        suffix = " [y/N]: "
    }

    for {
        fmt.Print(prompt + suffix)
        scanner.Scan()
        input := strings.ToLower(strings.TrimSpace(scanner.Text()))

        if input == "" {
            return defaultYes
        }

        if input == "y" || input == "yes" {
            return true
        }
        if input == "n" || input == "no" {
            return false
        }

        fmt.Println("Please answer 'y' or 'n'")
    }
}

// askList prompts for a list of strings
func askList(prompt string) []string {
    fmt.Println(prompt)
    fmt.Println("Enter items one per line. Enter blank line when done:")

    var items []string
    for {
        fmt.Print("> ")
        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())

        if input == "" {
            break
        }

        items = append(items, input)
    }

    return items
}

// askDayOfWeek prompts for a day of the week
func askDayOfWeek(prompt string) string {
    days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
    idx := askChoice(prompt, days, -1)
    return days[idx]
}

// askTime prompts for a time in HH:MM format
func askTime(prompt string) string {
    for {
        fmt.Print(prompt + " (HH:MM): ")
        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())

        parts := strings.Split(input, ":")
        if len(parts) != 2 {
            fmt.Println("Please enter time in HH:MM format (e.g., 18:30)")
            continue
        }

        hour, err1 := strconv.Atoi(parts[0])
        minute, err2 := strconv.Atoi(parts[1])

        if err1 != nil || err2 != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
            fmt.Println("Invalid time. Please enter HH:MM format (e.g., 18:30)")
            continue
        }

        return input
    }
}
