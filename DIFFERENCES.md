# Differences

Here are some of the noteworthy differences I have found between the Discord and Fluxer API so far.

## General

`nonce`:
- While Discord allows integer nonce values in some places (and presumably for compatibility you need to make sure this is preserved), on Fluxer only strings are allowed.


## REST API

`applications/@me`:
- The `owner` field is missing.
- The `oauth2/applications/@me` route is not present.
- The bio is a separate field in the `bot` object rather than the top level `description`. I haven't investigated how the top level `description` field is set but it remains null after setting a bio.

`users/{id}`:
- If `{id}` does not correspond to a real user it responds with DeletedUser#0000 with the provided ID unlike Discord which responds with the error message `Unknown User`.

## Gateway

Fluxer doesn't seem to support sharding apart from including some related constants and including `shards` as `1` in the gateway info endpoint for compatibility with Discord libs.

`READY` event:
- The `application` field is missing.
- The `resume_gateway_url` field is missing.
- The `shard` field is missing
