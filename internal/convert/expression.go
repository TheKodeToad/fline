package convert

import (
	"strings"

	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
	"github.com/disgoorg/snowflake/v2"
)

func EmojiToDiscord(emoji fluxer.Emoji) discord.Emoji {
	return discord.Emoji{
		ID: emoji.ID,
		Name: emoji.Name,
		Roles: []snowflake.ID{},
		RequireColons: misc.New(true),
		Managed: misc.New(false),
		Animated: emoji.Animated,
		Available: misc.New(true),
	}
}

func StickerToDiscord(sticker fluxer.Sticker) discord.Sticker {
	var formatType discord.StickerFormat
	// approximation
	if sticker.Animated {
		formatType = discord.StickerFormatGIF
	} else {
		formatType = discord.StickerFormatPNG
	}

	return discord.Sticker{
		ID: sticker.ID,
		Name: sticker.Name,
		Description: misc.New(sticker.Description),
		Tags: strings.Join(sticker.Tags, ","),
		Type: discord.StickerTypeGuild,
		FormatType: formatType,
		Available: misc.New(true),
	}
}
