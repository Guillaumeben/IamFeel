package main

import (
    "fmt"
    "os"
)

const (
    defaultUserConfigPath = "data/user_config.yaml"
    defaultDBPath         = "data/coach.db"
)

func main() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    command := os.Args[1]

    switch command {
    case "plan":
        if err := runPlanGeneration(); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
    case "help", "--help", "-h":
        printUsage()
        os.Exit(0)
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
        printUsage()
        os.Exit(1)
    }
}

func printUsage() {
    fmt.Println("IamFeel - AI-Powered Training Assistant")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("  iamfeel <command> [arguments]")
    fmt.Println()
    fmt.Println("Commands:")
    fmt.Println("  plan       Generate a training plan")
    fmt.Println("  help       Show this help message")
    fmt.Println()
    fmt.Println("Examples:")
    fmt.Println("  iamfeel plan")
}
