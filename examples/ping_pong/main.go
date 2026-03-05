package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
	"github.com/jurienhamaker/discordgoplus"
)

func init() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot load .env: %w", err))
	}
}

func main() {
	bot, err := discordgoplus.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	bot.Router.Register(&discordgoplus.Command{
		Name:        "ping_pong",
		Description: "Ping it!",
		Type:        discordgo.ChatApplicationCommand,
		Handler: discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
			_ = ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hi, I'm a bot built on Disgolf library.",
				},
			})
		}),
		MessageHandler: discordgoplus.MessageHandlerFunc(
			func(ctx *discordgoplus.MessageCtx) {
				_, _ = ctx.Reply("Hi, I'm a bot built on Disgolf library", true)
			},
		),

		Middlewares: []discordgoplus.Handler{
			discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
				fmt.Println("Middleware worked!")
				ctx.Next()
			}),
		},

		MessageMiddlewares: []discordgoplus.MessageHandler{
			discordgoplus.MessageHandlerFunc(
				func(ctx *discordgoplus.MessageCtx) {
					fmt.Println("Message niddleware worked!", ctx.Arguments)
					ctx.Next()
				},
			),
		},
	})
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(
		bot.Router.MakeMessageHandler(&discordgoplus.MessageHandlerConfig{
			Prefixes:      []string{"d.", "dis.", "discordgoplus."},
			MentionPrefix: true,
		}),
	)

	err = bot.Open()
	if err != nil {
		log.Fatal(fmt.Errorf("open exited with a error: %w", err))
	}
	defer bot.Close()
	err = bot.Router.Sync(bot.Session, "", "679281186975252480")
	if err != nil {
		log.Fatal(fmt.Errorf("cannot publish commands: %w", err))
	}
	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)
end:
	for {
		select {
		case <-stchan:
			break end
		default:
		}
		time.Sleep(time.Second)
	}
}
