package main

import (
    "fmt"
    "strings"
    "net/http"
    "net/http/cookiejar"
    "time"
    "net"

    "github.com/bwmarrin/discordgo"
)

var timeout = time.Duration(180 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, timeout)
}

func clonePodCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    if len(msgList) != 2 {
        s.ChannelMessageSend(m.ChannelID, "Usage is `!clone <pod name>`")
        return
    }

    jar, err := cookiejar.New(nil)
    if err != nil {
        fmt.Println("Error: ", err)
    }

    transport := &http.Transport{
        Dial: dialTimeout,
    }

    client := &http.Client{
        Jar: jar,
        Transport: transport,
    }

    postData := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, tomlConf.KaminoUser, tomlConf.KaminoPassword)
    loginUrl := "https://api.calpolyswift.org/login"

    req, err := http.NewRequest("POST", loginUrl, strings.NewReader(postData))
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error: ", err)
    }
    defer resp.Body.Close()

    cloneData := `{"template": "NCAE-Tryouts"}`
    cloneUrl := "https://api.calpolyswift.org/pod/clone/template"

    req, err = http.NewRequest("POST", cloneUrl, strings.NewReader(cloneData))
    req.Header.Set("Content-Type", "application/json")

    client.Jar.SetCookies(req.URL, resp.Cookies())

    resp, err = client.Do(req)
    if err != nil {
        fmt.Println("Error: ", err)
    }
    defer resp.Body.Close()

    s.ChannelMessageSend(m.ChannelID, "Cloning pod, check vSphere for progress.")
}
