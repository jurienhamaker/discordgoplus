package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jurienhamaker/discordgoplus"
)

type ExampleModule struct{}

func (ExampleModule) PingFunctional(s *discordgo.Session) time.Duration {
	return s.HeartbeatLatency()
}

func (m ExampleModule) Ping(ctx *discordgoplus.Ctx) {
	_ = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				":ping_pong: %v",
				ctx.HeartbeatLatency(),
			),
		},
	})
}

func (ExampleModule) PingMessage(ctx *discordgoplus.MessageCtx) {
	_, _ = ctx.Reply(
		fmt.Sprintf(
			":ping_pong: %v",
			ctx.HeartbeatLatency(),
		),
		false,
	)
}

func (m ExampleModule) Commands() []*discordgoplus.Command {
	return []*discordgoplus.Command{
		{
			Name:           "ping",
			Description:    "Get bot ping",
			Handler:        discordgoplus.HandlerFunc(m.Ping),
			MessageHandler: discordgoplus.MessageHandlerFunc(m.PingMessage),
		},
		{
			Name:        "ping_functional",
			Description: "Get bot ping",
			Handler: discordgoplus.HandlerFunc(func(ctx *discordgoplus.Ctx) {
				_ = ctx.Respond(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(
							":ping_pong: %v",
							m.PingFunctional(ctx.Session),
						),
					},
				})
			}),
			MessageHandler: discordgoplus.MessageHandlerFunc(
				func(ctx *discordgoplus.MessageCtx) {
					_, _ = ctx.Reply(
						fmt.Sprintf(
							":ping_pong: %v",
							m.PingFunctional(ctx.Session),
						),
						false,
					)
				},
			),
		},
	}
}
