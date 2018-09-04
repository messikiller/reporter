package main

import (
    "log"
    "time"
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

func isFileExists(filename string) bool {
    _, err := os.Stat(filename)

    if (err == nil || !os.IsNotExist(err)) {
        return true
    } else {
        return false
    }
}

func waitIfNotExists(filename string) {
    duration := 250 * time.Millisecond
    for {
		if isFileExists(filename) {
            return
        }

		select {
		case <-time.After(duration):
		}
	}
}

func watch(v configListItem, m configMail) {
    waitIfNotExists(v.Path)

    watcher, _ := fsnotify.NewWatcher()
    defer watcher.Close()

    watcher.Add(v.Path)

    f, _ := os.Open(v.Path)
    originStat, _ := os.Stat(v.Path)
    originSize := originStat.Size()
    f.Close()

    for {
        select {
        case ev := <-watcher.Events:
            if ev.Op & fsnotify.Remove == fsnotify.Remove {
                log.Println("File removed: ", ev.Name)

                watcher.Close()
                waitIfNotExists(ev.Name)
                watcher, _ = fsnotify.NewWatcher()
                watcher.Add(ev.Name)

                st, _ := os.Stat(v.Path)
                originSize = st.Size()

                if originSize > 0 {
                    buffer := make([]byte, originSize)

                    f, _ = os.Open(v.Path)
                    f.ReadAt(buffer, 0)
                    f.Close()

                    to := v.Emails
                    subject := v.Title
                    body := "<pre>" + string(buffer) + "</pre>"

                    go mail(m, to, subject, body)
                }
            }

            if ev.Op & fsnotify.Write == fsnotify.Write {
                stat, _ := os.Stat(v.Path)
                newSize := stat.Size()

                if newSize > originSize {
                    bufferSize := newSize - originSize
                    buffer := make([]byte, bufferSize)

                    f, _ = os.Open(v.Path)
                    f.ReadAt(buffer, originSize)
                    f.Close()

                    originSize = newSize

                    to := v.Emails
                    subject := v.Title
                    body := "<pre>" + string(buffer) + "</pre>"

                    go mail(m, to, subject, body)
                }
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
        go watch(v, r.Mail)
    }

    done := make(chan bool)
    <-done
}
