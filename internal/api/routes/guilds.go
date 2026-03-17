package apiroutes

import (
	"net/http"
	"net/url"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func guildsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/{guild_id}", api.ProxyHandler[any, fluxer.Guild]{
		Conf:   conf,
		Client: client,
		Path:   "/guilds/{guild_id}",
		MapResponse: func(guild fluxer.Guild) (any, error) {
			return convert.GuildToDiscord(guild), nil
		},
	})

	type banCreate struct {
		discord.GuildBanCreate
		auditLogReason string
	}

	router.Method("PUT", "/{guild_id}/bans/{user_id}", api.ProxyHandler[banCreate, api.NoContentResponse]{
		Client: client,
		Conf:   conf,
		Path:   "/guilds/{guild_id}/bans/{user_id}",
		DecodeRequest: func(req *http.Request) (banCreate, error) {
			create, err := api.DecodeRequestJSON[discord.GuildBanCreate](req)
			if err != nil {
				return banCreate{}, err
			}

			return banCreate{create, req.Header.Get("X-Audit-Log-Reason")}, nil
		},
		MapRequest: func(inCreate banCreate) (any, error) {
			outCreate := convert.GuildBanCreateToFluxer(inCreate.GuildBanCreate)

			if inCreate.auditLogReason != "" {
				// NOTE: fluxer has a separate message for the audit log and ban list
				// this replicates the discord behaviour of them being the same
				unescaped, err := url.PathUnescape(inCreate.auditLogReason)
				if err != nil {
					unescaped = inCreate.auditLogReason
				}

				outCreate.Reason = unescaped
			}

			return outCreate, nil
		},
		DecodeResponse: api.ExpectNoContentResponse,
	})

	memberRole := api.ProxyHandler[any, api.NoContentResponse]{
		Conf:           conf,
		Client:         client,
		Path:           "/guilds/{guild_id}/members/{user_id}/roles/{role_id}",
		DecodeResponse: api.ExpectNoContentResponse,
	}

	router.Method("PUT", "/{guild_id}/members/{user_id}/roles/{role_id}", memberRole)
	router.Method("DELETE", "/{guild_id}/members/{user_id}/roles/{role_id}", memberRole)

	return router
}
