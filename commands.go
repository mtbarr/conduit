package main

import (
	"conduit/i18n"

	"github.com/bwmarrin/discordgo"
)

type commandRegistrar interface {
	ApplicationCommandCreate(appID, guildID string, cmd *discordgo.ApplicationCommand, options ...discordgo.RequestOption) (*discordgo.ApplicationCommand, error)
	ApplicationCommandDelete(appID, guildID, cmdID string, options ...discordgo.RequestOption) error
}

type registeredCommand struct {
	cmd *discordgo.ApplicationCommand
	id  string
}

func registerReportBugCommand(registrar commandRegistrar, appID, guildID string) (*discordgo.ApplicationCommand, error) {
	return registrar.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
		Name:        i18n.T("reportbug_command_name"),
		Description: i18n.T("reportbug_command_desc"),
	})
}

func registerRequestFeatureCommand(registrar commandRegistrar, appID, guildID string) (*discordgo.ApplicationCommand, error) {
	return registrar.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
		Name:        i18n.T("requestfeature_command_name"),
		Description: i18n.T("requestfeature_command_desc"),
	})
}

func registerIssuesCommand(registrar commandRegistrar, appID, guildID string) (*discordgo.ApplicationCommand, error) {
	return registrar.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
		Name:        i18n.T("issues_command_name"),
		Description: i18n.T("issues_command_desc"),
	})
}

func deleteReportBugCommand(registrar commandRegistrar, appID, guildID, cmdID string) error {
	return registrar.ApplicationCommandDelete(appID, guildID, cmdID)
}

func deleteRequestFeatureCommand(registrar commandRegistrar, appID, guildID, cmdID string) error {
	return registrar.ApplicationCommandDelete(appID, guildID, cmdID)
}

func deleteIssuesCommand(registrar commandRegistrar, appID, guildID, cmdID string) error {
	return registrar.ApplicationCommandDelete(appID, guildID, cmdID)
}
