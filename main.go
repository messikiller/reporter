package main

import (
    "reporter/util"
    "fmt"
)

func main() {
    r := []util.ConfigItem{}
    util.ReadJson("./config.json", &r)

    for _, v := range r {
        fmt.Printf("title: %s, path: %s, emails: %s\n", v.Title, v.Path, v.Emails)
    }
}
