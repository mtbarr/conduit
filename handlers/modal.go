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

	customID := i.ModalSubmitData().CustomID
	log.Printf("Received modal submit: %q", customID)
	switch customID {
	case "modal_reportbug":
		handleReportBugModal(s, i)
	case "modal_requestfeature":
		handleRequestFeatureModal(s, i)
	default:
		log.Printf("Unknown modal: %q", customID)
	}
}

func handleReportBugModal(s interactionResponder, i *discordgo.InteractionCreate) {
	userID := getUserID(i)
	if userID != "" {
		allowed, remaining := checkAndMarkCooldown("reportbug_"+userID, reportbugCooldownDuration)
		if !allowed {
			respondEphemeral(s, i, i18n.Tf("reportbug_cooldown_message", formatCooldownRemaining(remaining)))
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

	_, err := createIssue(title, description, []string{"bug"})
	if err != nil {
		logError("create github issue", err)
		respondEphemeral(s, i, i18n.T("reportbug_issue_failed"))
		return
	}

	respondEphemeral(s, i, i18n.T("reportbug_issue_created_simple"))
}

func handleRequestFeatureModal(s interactionResponder, i *discordgo.InteractionCreate) {
	userID := getUserID(i)
	if userID != "" {
		allowed, remaining := checkAndMarkCooldown("requestfeature_"+userID, requestfeatureCooldownDuration)
		if !allowed {
			respondEphemeral(s, i, i18n.Tf("requestfeature_cooldown_message", formatCooldownRemaining(remaining)))
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
			case "feature_title":
				title = input.Value
			case "feature_description":
				description = input.Value
			}
		}
	}

	_, err := createIssue(title, description, []string{"feature-request"})
	if err != nil {
		logError("create github feature request issue", err)
		respondEphemeral(s, i, i18n.T("requestfeature_issue_failed"))
		return
	}

	respondEphemeral(s, i, i18n.T("requestfeature_issue_created_simple"))
}

func getUserID(i *discordgo.InteractionCreate) string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	if i.User != nil {
		return i.User.ID
	}
	return ""
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

func respondEphemeralIssues(s interactionResponder, i *discordgo.InteractionCreate) {
	log.Println("Fetching issues from GitHub...")
	issues, err := github.ListIssues(10)
	if err != nil {
		logError("list github issues", err)
		respondEphemeral(s, i, i18n.T("issues_failed"))
		return
	}

	if len(issues) == 0 {
		respondEphemeral(s, i, i18n.T("issues_no_issues"))
		return
	}

	content := i18n.T("issues_header") + "\n\n"
	for _, issue := range issues {
		labels := ""
		if len(issue.Labels) > 0 {
			labels = issue.Labels[0]
		}
		content += i18n.Tf("issue_format", issue.Number, issue.Title, labels) + "\n"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logError("send issues ephemeral response", err)
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
