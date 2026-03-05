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
		Name:        "subcommands",
		Description: "Lo and behold, subcommands are coming!",
		Type:        discordgo.ChatApplicationCommand,
		MessageMiddlewares: []discordgoplus.MessageHandler{
			discordgoplus.MessageHandlerFunc(func(ctx *discordgoplus.MessageCtx) {
				fmt.Println("middleware")
				ctx.Next()
			}),
		},
		Middlewares: []discordgoplus.Handler{
			discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
				fmt.Println("middleware")
				ctx.Next()
			}),
		},
		SubCommands: discordgoplus.NewRouter([]*discordgoplus.Command{
			{
				Name:        "group",
				Description: "Subcommand group",
				MessageMiddlewares: []discordgoplus.MessageHandler{
					discordgoplus.MessageHandlerFunc(func(ctx *discordgoplus.MessageCtx) {
						fmt.Println("group middleware")
						ctx.Next()
					}),
				},
				Middlewares: []discordgoplus.Handler{
					discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
						fmt.Println("group middleware")
						ctx.Next()
					}),
				},
				SubCommands: discordgoplus.NewRouter([]*discordgoplus.Command{
					{
						Name:        "subcommand",
						Description: "Subcommand in a subcommand group",
						Handler: discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
							_ = ctx.Respond(&discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "hi (group)",
								},
							})
						}),
						MessageHandler: discordgoplus.MessageHandlerFunc(
							func(ctx *discordgoplus.MessageCtx) {
								_, _ = ctx.Reply("hi (group)", false)
							},
						),
						MessageMiddlewares: []discordgoplus.MessageHandler{
							discordgoplus.MessageHandlerFunc(
								func(ctx *discordgoplus.MessageCtx) {
									fmt.Println("individual middleware")
									ctx.Next()
								},
							),
						},
						Middlewares: []discordgoplus.Handler{
							discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
								fmt.Println("individual middleware")
								ctx.Next()
							}),
						},
					},
				}),
				MessageHandler: discordgoplus.MessageHandlerFunc(
					func(ctx *discordgoplus.MessageCtx) {
						_, _ = ctx.Reply("hi (group default)", false)
					},
				),
			},
			{
				Name:        "subcommand",
				Description: "Just a subcommand",
				Handler: discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
					_ = ctx.Respond(&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{Content: "hi"},
					})
				}),
				MessageHandler: discordgoplus.MessageHandlerFunc(
					func(ctx *discordgoplus.MessageCtx) {
						_, _ = ctx.Reply("hi", false)
					},
				),
				MessageMiddlewares: []discordgoplus.MessageHandler{
					discordgoplus.MessageHandlerFunc(func(ctx *discordgoplus.MessageCtx) {
						fmt.Println("individual middleware (2nd level)")
						ctx.Next()
					}),
				},
				Middlewares: []discordgoplus.Handler{
					discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
						fmt.Println("individual middleware (2nd level)")
						ctx.Next()
					}),
				},
			},
		}),
		MessageHandler: discordgoplus.MessageHandlerFunc(
			func(ctx *discordgoplus.MessageCtx) {
				_, _ = ctx.Reply("hi (default)", false)
			},
		),
	})
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(bot.Router.MakeMessageHandler(&discordgoplus.MessageHandlerConfig{
		Prefixes:      []string{"d.", "dis.", "discordgoplus."},
		MentionPrefix: true,
	}))

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
