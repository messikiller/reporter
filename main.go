package main

import (
    "reporter/util"
    "github.com/hpcloud/tail"
)

func watch(filename string, to string, subject string) {
    tails, err := tail.TailFile(filename, tail.Config{
        ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
    })

    if err != nil {
        panic(err)
        return
    }

    var line *tail.Line
    var body string

    for true {
        select {
        case <-tails.Lines:
            line = <-tails.Lines
            body += line.Text
        default:
            if body != "" {
                go util.Mail(to, subject, body)
                body = ""
            }
        }
    }


}

func main() {
    r := []util.ConfigItem{}
    util.ReadJson("./config.json", &r)

    for _, v := range r {
        go watch(v.Path, v.Emails, "[" + v.Title + "]" + " Exception Found")
    }

    ch := make(chan bool)
    <-ch
}
