package util

import (
    "io/ioutil"
    "encoding/json"
)

type ConfigItem struct {
    Title string
    Path string
    Emails string
}

func ReadJson(path string, r *[]ConfigItem) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }

    err = json.Unmarshal(data, &r)
    if err != nil {
        panic(err)
    }
}
