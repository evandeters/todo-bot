package main

import (
    "strings"

    "github.com/bwmarrin/discordgo"
)

var (
    activeCommands = make(map[string]command)
    disabledCommands = make(map[string]command)
)

type command struct {
    Name string
    Help string

    Exec func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
    msgList := strings.Fields(message)
    if len(msgList) == 0 {
        return
    }

    commandName := strings.ToLower(func() string {
        if strings.HasPrefix(message, " ") {
            return " " + msgList[0]
        }
        return msgList[0]
    }())

    if command, ok := activeCommands[commandName]; ok && commandName == strings.ToLower(command.Name) {
        command.Exec(s, m, msgList)
        return
    }
}

func (c command) add() command {
    activeCommands[strings.ToLower(c.Name)] = c
    return c
}

func newCommand(name string, f func(*discordgo.Session, *discordgo.MessageCreate, []string)) command {
    return command{
        Name: name,
        Exec: f,
    }
} 

func (c command) setHelp(help string) command {
    c.Help = help
    return c
}

func getCommand(name string) *command {
    if c, ok := activeCommands[name]; ok {
        return &c
    }
    return nil
}
