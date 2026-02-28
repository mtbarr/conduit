package main

import (
	"conduit/handlers"
	"conduit/i18n"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: could not load .env file, falling back to environment variables")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	session.Identify.Intents = discordgo.IntentsGuilds

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %s#%s", r.User.Username, r.User.Discriminator)
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			handlers.HandleCommand(s, i)
		case discordgo.InteractionModalSubmit:
			handlers.HandleModal(s, i)
		}
	})

	if err := session.Open(); err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("Error closing Discord session: %v", err)
		}
	}()

	cmd, err := registerReportBugCommand(session, session.State.User.ID, "")
	if err != nil {
		log.Fatalf("Error registering /reportbug command: %v", err)
	}
	log.Printf("Language: %s", i18n.CurrentLanguage())
	log.Println("Bot is running. Press Ctrl+C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Shutting down...")
	if err := deleteReportBugCommand(session, session.State.User.ID, "", cmd.ID); err != nil {
		log.Printf("Error removing /reportbug command: %v", err)
	}
	log.Println("Bot stopped.")
}
