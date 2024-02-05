package main

import (
    "fmt"
    "strings"
    "strconv"

    "github.com/Clinet/discordgo-embed"
    "github.com/bwmarrin/discordgo"
)

func addTodoCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    if len(msgList) < 1 {
        s.ChannelMessageSend(m.ChannelID, "User and todo not provided")
        return
    }

    user := msgList[1]
    todo := strings.Join(msgList[2:], " ")
    
    err := addTodo(user, todo)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error adding todo")
        return
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Todo Added")
    embed.SetColor(0x38fcec)
    embed.AddField("User", user)
    embed.AddField("Task", todo)

    s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func getAllTodosCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    todos, err := getAllTodos()
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error getting todos")
        return
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Todos")
    embed.SetColor(0x38fcec)

    for _, t := range todos {
        embed.AddField(fmt.Sprintf("Id: %d", t.Id), fmt.Sprintf("User: %s\n Task: %s\n Completed: %t", t.User, t.Task, t.Completed))
    }

    s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func getTodosByUserCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    if len(msgList) < 2 {
        s.ChannelMessageSend(m.ChannelID, "User not provided")
        return
    }

    user := msgList[1]
    
    todos, err := getTodos(user)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error getting todos")
        return
    }

    if len(todos) == 0 {
        s.ChannelMessageSend(m.ChannelID, "No todos found for " + user)
        return
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Todos for " + user)
    embed.SetColor(0x38fcec)

    for _, t := range todos {
        embed.AddField(fmt.Sprintf("Id: %d", t.Id), fmt.Sprintf("Task: %s\n Completed: %t", t.Task, t.Completed))
    }

    s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func removeTodoByUserCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    if len(msgList) < 2 {
        s.ChannelMessageSend(m.ChannelID, "User not provided")
        return
    }

    user := msgList[1]
    
    err := deleteTodos(user)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error removing todo")
        return
    }

    s.ChannelMessageSend(m.ChannelID, "Todos removed for " + user)
}

func completeTodoCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    id, err := strconv.Atoi(msgList[1])
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Invalid id")
        return
    }
    
    err = completeTodoById(id)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error completing todo")
        return
    }
    s.ChannelMessageSend(m.ChannelID, "Todo completed")
}

func updateTodoCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    id, err := strconv.Atoi(msgList[1])
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Invalid id")
        return
    }
    todo := strings.Join(msgList[2:], " ")
    
    err = updateTodoById(id, todo)
    if err != nil {
        s.ChannelMessageSend(m.ChannelID, "Error updating todo")
        return
    }
    s.ChannelMessageSend(m.ChannelID, "Todo updated")
}
