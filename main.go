package main

import (
        "fmt"
        "log"
        "os"
        "os/signal"
        "syscall"
        "time"
)

func main() {
        log.Println("Hello World Service started!")

        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

        go func() {
                <-sigChan
                log.Println("Hello World Service stopping...")
                os.Exit(0)
        }()

        for {
                fmt.Println("Hello from Hello World Service!")
                time.Sleep(5 * time.Second)
        }
}
