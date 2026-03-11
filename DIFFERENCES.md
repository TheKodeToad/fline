# Differences

Here are some of the noteworthy differences I have found between the Discord and Fluxer API so far.

## REST API

`applications/@me`:
- The `owner` field is missing.
- The `oauth2/applications/@me` route is not present.
- The bio is a separate field in the `bot` object rather than the top level `description`. I haven't investigated how the top level `description` field is set but it remains null after setting a bio.

# Gateway

Fluxer doesn't seem to support sharding apart from including some related constants and including `shards` as `1` in the gateway info endpoint for compatibility with Discord libs.

`READY` event:
- The `application` field is missing.
- The `resume_gateway_url` field is missing.
- The `shard` field is missing
