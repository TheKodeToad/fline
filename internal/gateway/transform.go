package gateway

import (
	"fmt"

	"github.com/TheKodeToad/fine/internal/discord"
)

// packetToDiscord converts from a Fluxer packet to a Discord packet.
func packetToDiscord(packet discord.Packet) (discord.Packet, bool) {
	fmt.Printf("converting packet to discord %+v %s\n", packet, string(packet.Data))

	switch packet.Opcode {
	case discord.GatewayOpHello:
		return packet, true
	}

	return discord.Packet{}, false
}

// packetToFluxer converts from a Discord packet to a Fluxer packet.
func packetToFluxer(packet discord.Packet) (discord.Packet, bool) {
	fmt.Printf("converting packet to fluxer %+v %s\n", packet, string(packet.Data))

	return discord.Packet{}, false
}
