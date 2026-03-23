# Fline

A <ins>highly experimental</ins>, work-in-progress API adapter specifically for running Discord bots on Fluxer.

The current goal is to allow Discord bots which only use features available on Fluxer to work smoothly without many changes other than modifying a few URLs in the code.

You may experience bugs! You will experience bugs! In fact, if there aren't any bugs, that's a bug! Bear in mind that I am just a nerd making this to learn cool stuff and my previous experience with the Fluxer and Discord API is limited.

## Isn't Fluxer's API already compatible with Discord??

There certainly are more simularities than differences, but the small differences which break assumptions made by Discord bots and libraries add up. From my testing most libraries can't make it far without throwing an error when connecting directly to the Fluxer API.

Often there are some required properties which libraries reasonably assume exist which Fluxer doesn't provide - and a sensible value can be inferred from context or a placeholder will work well enough.
There's also properties which are documented to be optional but libraries assume are present - and even properties which are documented to be required but can actually be omitted in the request body.

Many of these differences I am documenting [here](./DIFFERENCES.md) in case Fluxer would like to "fix" their API to be the same in some areas (how much of this is intentional I am not sure).

## Shouldn't you just port bots to use the Fluxer API properly?

That would probably produce better results, but the point of this project is to allow a number of bots to work without much effort.
