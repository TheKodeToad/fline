package fline

const FluxerAPIVersion = "1"

const GatewayPath = "/gateway"

func GatewayURL(host string) string {
	return "ws://" + host + GatewayPath
}
