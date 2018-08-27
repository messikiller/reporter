package main

import (
    "os"
    "flag"
    "io/ioutil"
    "encoding/json"
    "crypto/tls"

    "github.com/fsnotify/fsnotify"
    gomail "github.com/go-mail/gomail"
)

type configMail struct {
    From string
    Host string
    Port int
    Username string
    Password string
}

type configListItem struct {
    Title string
    Path string
    Emails string
}

type config struct {
    Mail configMail
    List []configListItem
}

func readJson(path string, r *config) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }

    err = json.Unmarshal(data, &r)
    if err != nil {
        panic(err)
    }
}

func mail(mail configMail, to string, subject string, body string) {
    from := mail.From
    host := mail.Host
    port := mail.Port
    username := mail.Username
    password := mail.Password

    m := gomail.NewMessage()

    m.SetHeader("From", from)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer(host, port, username, password)
    d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
}

func watch(filename string, to string, subject string, m configMail) {
    watcher, _ := fsnotify.NewWatcher();
    defer watcher.Close();

    watcher.Add(filename);

    f, _ := os.Open(filename)
    size, _ := f.Seek(0, os.SEEK_END)

    for {
        select {
        case ev := <-watcher.Events:
            if ev.Op & fsnotify.Write == fsnotify.Write {
                stat, _ := os.Stat(filename)
                newSize := stat.Size()

                bufferSize := newSize - size
                buffer := make([]byte, bufferSize)
                f.ReadAt(buffer, size)

                size = newSize

                body := "<pre>" + string(buffer) + "</pre>"

                go mail(m, to, subject, body)
            }
        }
    }
}

func main() {
    configPath := flag.String("c", "/etc/reporter.json", "config.json file path")

    flag.Parse()

    r := config{}
    readJson(*configPath, &r)

    for _, v := range r.List {
        go watch(v.Path, v.Emails, "【" + v.Title + "】" + " Exception Found", r.Mail)
    }

    ch := make(chan bool)
    <-ch
}
