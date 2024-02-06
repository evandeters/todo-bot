package main

import (
    "fmt"
    "strings"
    "net/http"
    "net/http/cookiejar"
    "time"
    "net"
    "encoding/json"
    "io"

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

    client, resp := login()
    pods, err := getPods(client)
    if err != nil {
        fmt.Println("Error: ", err)
    }

    podName := ""
    for _, p := range pods {
        if p == msgList[1] {
            podName = p.(string)
            break
        }
    }
    if podName == "" {
        s.ChannelMessageSend(m.ChannelID, "Pod not found.")
        return
    }

    cloneData := fmt.Sprintf(`{"template": "%s"}`, podName)
    cloneUrl := "https://api.calpolyswift.org/pod/clone/template"

    s.ChannelMessageSend(m.ChannelID, "Cloning pod, check vSphere for progress.")
    req, err := http.NewRequest("POST", cloneUrl, strings.NewReader(cloneData))
    req.Header.Set("Content-Type", "application/json")

    client.Jar.SetCookies(req.URL, resp.Cookies())

    resp, err = client.Do(req)
    if err != nil {
        fmt.Println("Error: ", err)
    }
    defer resp.Body.Close()
}

func getPodsCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    if len(msgList) != 1 {
        s.ChannelMessageSend(m.ChannelID, "Usage is `!pods`")
        return
    }

    client, _ := login()
    pods, err := getPods(client)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error getting pods.")
        fmt.Println("Error: ", err)
        return
    }

    podList := "**Pods**:\n"
    for _, p := range pods {
        podList += fmt.Sprintf("- %s\n", p)
    }

    s.ChannelMessageSend(m.ChannelID, podList)
}

func getPods(c *http.Client) ([]interface{}, error) {
    url := "https://api.calpolyswift.org/view/templates/preset"

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.Do(req)
    if err != nil {
        return nil, err
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    jsonData := make(map[string]interface{})
    err = json.Unmarshal(body, &jsonData)
    if err != nil {
        return nil, err
    }

    pods := jsonData["templates"].([]interface{})

    return pods, nil

}

func login() (*http.Client, *http.Response) {
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

    return client, resp
}
