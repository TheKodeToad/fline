# Differences

Here are some of the noteworthy differences I have found between the Discord and Fluxer API so far.

## General

`nonce`:
- While Discord allows integer nonce values in some places (and presumably for compatibility you need to make sure this is preserved), on Fluxer only strings are allowed.

Roles:
- There are no role tags.
- There is no `colors` field - which makes sense without role gradients - but Discord has deprecated the color field and made this field required in the response.

Permissions:
- The permission `UPDATE_RTC_REGION` (`1 << 53` - "Set Voice Region") is not present on Discord. If Discord adds another permission it could conflict.

## REST API

`allowed_mentions`:
- The mere precense of the object does not disable any mentions, unlike Discord where `{}` disables (to my knowledge) all mentions.
- `parse` defaults to ["users", "roles", "everyone"] unless either the `roles` or `users` key exists (or both). `replied_user` defaults to `true` and must be explicitly set to `false`.

`X-Audit-Log-Reason`:
- Unlike Discord it is not URL encoded.

GET `applications/@me`:
- The `owner` field is missing.
- The `oauth2/applications/@me` route is not present.
- The bio is a separate field in the `bot` object rather than the top level `description`. I haven't investigated how the top level `description` field is set but it remains null after setting a bio.

POST `channels/{channel_id}/messages`
- On Discord, omitting the required property on `footer`, `image`, `thumbnail` or `author` does not yield an error response - the property with the missing field is ignored instead. Fluxer does not replicate this.
- With a multipart body, only `payload_json` and `files[n]` is supported - and the former is required. While many (most?) implementations only use these, Discord also allows the rest of the top-level non-object properties to be set. You can even have multiple `sticker_ids` by specifying it multiple times. Discord also doesn't require attachment metadata, while Fluxer does, requiring you to specify the filenames again. Fluxer also infers the content type based on the filename with no way to specify it explicitly.

GET `guilds/{guild_id}`
- `roles`, `emojis` and `stickers` are missing

POST `guilds/{guild_id}/bans/{user_id}`
- The audit log reason (`X-Audit-Log-Reason` header) is separate from the actual reason stored in the entry (`reason` property).

GET `users/{user_id}`:
- If `{user_id}` does not correspond to a real user it responds with DeletedUser#0000 with the provided ID unlike Discord which responds with the error message `Unknown User`.

## Gateway

Fluxer doesn't seem to support sharding apart from including some related constants and including `shards` as `1` in the gateway info endpoint for compatibility with Discord libs.

`READY` event:
- The `application` field is missing.
- The `resume_gateway_url` field is missing.
- The `shard` field is missing

`GUILD_CREATE` event:
- Most of the guild's properties are within `properties`.
