> [!IMPORTANT]
> I have moved on to working a [Go library for Fluxer from scratch](https://github.com/fluxer-flo/flo) (should be available at some point), and this project seems like it's probably higher effort and has worse results than porting bots so I will likely not continue

# Fline

A <ins>highly experimental</ins>, work-in-progress API adapter (with some possible glaring issues) specifically for running Discord bots on Fluxer.

The current goal is to allow Discord bots which only use features also available on Fluxer for the most part to work smoothly without many changes other than modifying a few URLs in the code.

While I have not yet reached this goal, I am happy to have got some functionality from [Zeppelin](https://zeppelin.gg/) and its dashboard working with a few URL changes and deleted code here and there.

<img width="602" height="751" alt="image" src="https://github.com/user-attachments/assets/d69774e5-a3c7-411b-b882-acf6db1ae86e" />

You may experience bugs! You will experience bugs! In fact, if there aren't any bugs, that's a bug! Bear in mind that I am just a nerd making this to learn cool stuff and my previous experience with the Fluxer and Discord API is limited.

Right now this project is heavily limited by Fluxer lacking interactions.

## Usage

You will need to install [a toolchain for the Go programming language](https://go.dev/dl/).

First run `go build` to build and then it can be run with `./fline`. This will start a HTTP/websocket server running on http://localhost:8080 - the API at /api and gateway at /gateway.

Environment variables can be provided through the command line or `.env`; you can see which are available in [.env.example](./.env.example). Most of the defaults should already be good enough if you want to use the official instance of Fluxer, but you may want to change the port that is being listened on with `FLINE_LISTEN_ADDR=:1234`.

Next, you need to get the bot you wish to run to connect to Fline instead of the real Discord API. With Discord.js this is simple:
```js
const client = new Client({
    // ...
    rest: {
        api: "http://localhost:<enter port here>/api",
    }
});
```
Unfortunately not all libraries have a way to override the URL so you may need to modify the library code or some other shenanigans.

## Why?

Isn't Fluxer's API already compatible with Discord? Well, it is very similar which is nice for a project like this - but in practice [the many differences](./DIFFERENCES.md) cause things not to work right - especially missing fields.

Porting bots to Fluxer (and taking advantage of Fluxer specific features, more of which will likely come) is probably better, but this project is interesting to me and may be useful to others.
