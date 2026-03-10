package fine

const GatewayPath = "/gateway"

func GatewayURL(host string) string {
	return "ws://" + host + GatewayPath
}
