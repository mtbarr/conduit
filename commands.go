package main

import (
	"conduit/i18n"

	"github.com/bwmarrin/discordgo"
)

type commandRegistrar interface {
	ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error)
	ApplicationCommandDelete(appID, guildID, cmdID string, options ...discordgo.RequestOption) error
}

func registerReportBugCommand(registrar commandRegistrar, appID, guildID string) (*discordgo.ApplicationCommand, error) {
	return registrar.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
		Name:        i18n.T("command_name"),
		Description: "Report a bug — opens a form to create a GitHub issue",
	})
}

func deleteReportBugCommand(registrar commandRegistrar, appID, guildID, cmdID string) error {
	return registrar.ApplicationCommandDelete(appID, guildID, cmdID)
}
