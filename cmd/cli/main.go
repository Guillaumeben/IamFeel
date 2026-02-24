package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    command := os.Args[1]

    switch command {
    case "onboard":
        if err := runOnboarding(); err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
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
    fmt.Println("  onboard    Run the onboarding wizard to set up your profile")
    fmt.Println("  plan       Generate a training plan")
    fmt.Println("  help       Show this help message")
    fmt.Println()
    fmt.Println("Examples:")
    fmt.Println("  iamfeel onboard")
    fmt.Println("  iamfeel plan")
}
