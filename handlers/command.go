package handlers

import (
	"conduit/i18n"

	"github.com/bwmarrin/discordgo"
)

func HandleCommand(s interactionResponder, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if i.ApplicationCommandData().Name != i18n.T("command_name") {
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal_reportbug",
			Title:    i18n.T("modal_title"),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "bug_title",
							Label:       i18n.T("modal_title_label"),
							Style:       discordgo.TextInputShort,
							Placeholder: i18n.T("modal_title_placeholder"),
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
							Label:       i18n.T("modal_desc_label"),
							Style:       discordgo.TextInputParagraph,
							Placeholder: i18n.T("modal_desc_placeholder"),
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
