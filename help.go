package main

import (
    "github.com/Clinet/discordgo-embed"
    "github.com/bwmarrin/discordgo"
)

func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
    embed := embed.NewEmbed()
    embed.SetTitle("Command Help")
    embed.SetColor(0x38fcec)

    for _, c := range activeCommands {
        embed.AddField(c.Name, c.Help)
    }
    s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
