package main

import (
    "fmt"
    "time"
)

func main() {
    currentTime := time.Now()

    fmt.Println("(" + currentTime.Format("2006-01-02") + ")")
}
