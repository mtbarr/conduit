package handlers

import (
	"log"
	"time"

	"conduit/github"
	"conduit/i18n"

	"github.com/bwmarrin/discordgo"
)

var createIssue = github.CreateIssue

func HandleModal(s interactionResponder, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}
	if i.ModalSubmitData().CustomID != "modal_reportbug" {
		return
	}

	userID := ""
	if i.Member != nil && i.Member.User != nil {
		userID = i.Member.User.ID
	} else if i.User != nil {
		userID = i.User.ID
	}
	if userID != "" {
		allowed, remaining := checkAndMarkCooldown(userID)
		if !allowed {
			respondEphemeral(s, i, i18n.Tf("cooldown_message", formatCooldownRemaining(remaining)))
			return
		}
	}

	var title, description string
	for _, comp := range i.ModalSubmitData().Components {
		row, ok := comp.(*discordgo.ActionsRow)
		if !ok {
			continue
		}
		for _, rowComp := range row.Components {
			input, ok := rowComp.(*discordgo.TextInput)
			if !ok {
				continue
			}
			switch input.CustomID {
			case "bug_title":
				title = input.Value
			case "bug_description":
				description = input.Value
			}
		}
	}

	issueURL, err := createIssue(title, description)
	if err != nil {
		logError("create github issue", err)
		respondEphemeral(s, i, i18n.T("issue_failed"))
		return
	}

	_ = issueURL
	respondEphemeral(s, i, i18n.T("issue_created_simple"))
}

func respondEphemeral(s interactionResponder, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logError("send ephemeral response", err)
		return
	}
	if ephemeralDeleteDelay > 0 {
		go func() {
			time.Sleep(ephemeralDeleteDelay)
			_ = s.InteractionResponseDelete(i.Interaction)
		}()
	}
}

func logError(context string, err error) {
	log.Printf("ERROR [%s]: %v", context, err)
}
