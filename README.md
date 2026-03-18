# Fline

An experimental, work-in-progress and moderately serious API compatibility layer specifically for running Discord bots on Fluxer. The name stands for "Fline is not an emulator" (this might seem familiar).

The goal is currently to allow bots which use the subset of features both on Fluxer and Discord to work smoothly with few code changes!
I will also likely investigate the viability of supporting some of the features of the interactions API (perhaps slash commands could be processed from regular messages beginning with `/`).

## Isn't Fluxer's API already compatible with Discord??

There certainly are more simularities than differences, but the minor differences which break assumptions made by Discord bots and libraries add up. From my testing most libraries can't make it far without crashing when connecting directly to the Fluxer API.

Often there are some required properties which libraries reasonably assume exist which Fluxer doesn't provide - and a sensible value can be inferred from context or a placeholder will work well enough.
There's also properties which are documented to be optional but libraries assume are present, and even properties which are documented to be required but can actually be omitted in the request body.

Many of these differences I am documenting [here](./DIFFERENCES.md) in case Fluxer would like to "fix" their API to be the same in some areas (how much of this is intentional I am not sure).

## Shouldn't you just port bots to use the Fluxer API properly?

That would probably produce better results, but the point of this project is to allow a number of bots to work without much effort.
