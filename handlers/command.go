package handlers

import (
	"log"

	"conduit/i18n"

	"github.com/bwmarrin/discordgo"
)

func HandleCommand(s interactionResponder, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	log.Printf("Received command: %q", commandName)
	switch commandName {
	case i18n.T("reportbug_command_name"):
		handleReportBugCommand(s, i)
	case i18n.T("requestfeature_command_name"):
		handleRequestFeatureCommand(s, i)
	case i18n.T("issues_command_name"):
		handleIssuesCommand(s, i)
	default:
		log.Printf("Unknown command: %q (expected one of: %q, %q, %q)",
			commandName,
			i18n.T("reportbug_command_name"),
			i18n.T("requestfeature_command_name"),
			i18n.T("issues_command_name"))
	}
}

func handleReportBugCommand(s interactionResponder, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal_reportbug",
			Title:    i18n.T("reportbug_modal_title"),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "bug_title",
							Label:       i18n.T("reportbug_modal_title_label"),
							Style:       discordgo.TextInputShort,
							Placeholder: i18n.T("reportbug_modal_title_placeholder"),
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
							Label:       i18n.T("reportbug_modal_desc_label"),
							Style:       discordgo.TextInputParagraph,
							Placeholder: i18n.T("reportbug_modal_desc_placeholder"),
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
		logError("respond with reportbug modal", err)
	}
}

func handleRequestFeatureCommand(s interactionResponder, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modal_requestfeature",
			Title:    i18n.T("requestfeature_modal_title"),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "feature_title",
							Label:       i18n.T("requestfeature_modal_title_label"),
							Style:       discordgo.TextInputShort,
							Placeholder: i18n.T("requestfeature_modal_title_placeholder"),
							Required:    true,
							MinLength:   1,
							MaxLength:   256,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "feature_description",
							Label:       i18n.T("requestfeature_modal_desc_label"),
							Style:       discordgo.TextInputParagraph,
							Placeholder: i18n.T("requestfeature_modal_desc_placeholder"),
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
		logError("respond with requestfeature modal", err)
	}
}

func handleIssuesCommand(s interactionResponder, i *discordgo.InteractionCreate) {
	respondEphemeralIssues(s, i)
}
