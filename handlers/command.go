package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func HandleCommand(s interactionResponder, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name != "reportbug" {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal_reportbug",
			Title:    "Report a Bug",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "bug_title",
							Label:       "Title",
							Style:       discordgo.TextInputShort,
							Placeholder: "Short summary of the bug",
							Required:    true,
							MinLength:   1,
							MaxLength:   256,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "bug_description",
							Label:       "Description",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Detailed description of the bug",
							Required:    true,
							MinLength:   1,
							MaxLength:   4000,
						},
					},
				},
			},
		},
	})
	if err != nil {
		logError("respond with modal", err)
	}
}
