package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
)

func WebhookToDiscord(webhook fluxer.Webhook) discord.Webhook {
	return discord.Webhook{
		ID:        webhook.ID,
		Type:      discord.WebhookTypeIncoming,
		GuildID:   &webhook.ChannelID,
		ChannelID: &webhook.ChannelID,
		User:      misc.New(UserPartialToDiscord(webhook.User)),
		Name:      &webhook.Name,
		Avatar:    webhook.Avatar,
		Token:     &webhook.Token,
	}
}

func WebhookExecuteToDiscord(exec discord.WebhookExecute) (discord.WebhookExecute, bool) {
	if exec.AllowedMentions != nil {
		exec.AllowedMentions = misc.New(AllowedMentionsToFluxer(*exec.AllowedMentions))
	}

	attachments, files, ok := UploadsToFluxer(exec.Attachments, exec.Files)
	if !ok {
		return discord.WebhookExecute{}, false
	}

	exec.Attachments = attachments
	exec.Files = files

	return exec, true
}
