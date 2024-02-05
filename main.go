package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
    "strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
    prefix = "!"
    db = ConnectDB()
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

    newCommand("todo", addTodoCommand).setHelp("Add a new todo. Usage is '!todo <user> <task>'").add()
    newCommand("get-todo", getAllTodosCommand).setHelp("Get all todos.").add()
    newCommand("get-user-todo", getTodosByUserCommand).setHelp("Get all todos by user. Usage is '!get-user-todo <user>'").add()
    newCommand("remove-user-todo", removeTodoByUserCommand).setHelp("Remove all todos by user. Usage is '!remove-user-todo <user>'").add()
    newCommand("complete-todo", completeTodoCommand).setHelp("Complete a todo. Usage is '!complete-todo <id>'").add()
    newCommand("help", helpCommand).setHelp("Get help.").add()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

    err = db.AutoMigrate(&todo{})
    if err != nil {
        fmt.Println("Error migrating database:", err)
    }


	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.Bot {
        return
    }

    if !strings.HasPrefix(m.Content, prefix) {
        return
    }

    parseCommand(s, m, m.Content[len(prefix):])
}

