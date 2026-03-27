# Fline

A <ins>highly experimental</ins>, work-in-progress API adapter specifically for running Discord bots on Fluxer.

The current goal is to allow Discord bots which only use features also available on Fluxer for the most part to work smoothly without many changes other than modifying a few URLs in the code.

You may experience bugs! You will experience bugs! In fact, if there aren't any bugs, that's a bug! Bear in mind that I am just a nerd making this to learn cool stuff and my previous experience with the Fluxer and Discord API is limited.

Right now this project is heavily limited by Fluxer lacking interactions.

## Usage

First run `go build` to build and then it can be run with `./fline`.

Environment variables can be provided through the command line or `.env`; you can see which are available in [.env.example](./.env.example). Most of the defaults should already be good enough if you want to use the official instance of Fluxer, but you may want to change the port that is being listened on with `FLINE_LISTEN_ADDR=:1234` (it defaults to 8080).

## Why?

Isn't Fluxer's API already compatible with Discord? Well, it is very similar which is nice for a project like this - but in practice the many differences cause things not to work right - especially missing fields.

Porting bots to Fluxer (and taking advantage of Fluxer specific features, more of which will likely come) is probably better, but this project is interesting to me and may be useful to others.
