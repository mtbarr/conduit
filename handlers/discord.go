package handlers

import "github.com/bwmarrin/discordgo"

type interactionResponder interface {
	InteractionRespond(*discordgo.Interaction, *discordgo.InteractionResponse, ...discordgo.RequestOption) error
	InteractionResponseDelete(*discordgo.Interaction, ...discordgo.RequestOption) error
}
