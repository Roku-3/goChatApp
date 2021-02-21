package main
import (
    "time"
)

// messageは１つのメッセージを表す。
type message struct {
    Name    string
    Message string
    When    time.Time
}

