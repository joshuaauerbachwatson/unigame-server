# unigame-server

This repo contains the source for a game playing backend that can serve a wide variety of possible games.  The server is set up to be deployed on DigitalOcean App Platform but could be run in other ways.

The server supports
- simple authorization checks using auth0.  In order to connect, a JWT containing a valid auth0 access token must be presented.  To obtain this token, users must go through an auth0 login.
- communication via websocket once the authorization check has passed
- an indefinite number of ongoing games, each game being identified by a game token.  Players must agree on the game token by means outside the server (the server has no "social" functions).  The game token is separate from the authorization mechanism (anyone who passes the general authentication and knows the game token can join the game)
- maintenance of a "number of players" per game; the game starts once that many players have joined
- maintenance of a list of players for each game
- multicasting a simple text chat amongst the players, which commences even before the game is started
- multicasting of "game states" amongst the players.  The game states are not interpreted by the server, allowing almost any multiplayer game to be supported.
- a keepalive mechanism to detect lost players
- until all players have joined, a garbage collection mechanism that will delete incomplete games

Games capable of being played by `unigame-server` are supported by a Swift (iOS and Mac) app framework called [`unigame`](https://github.com/joshuaauerbachwatson/unigame).

Games currently built on the `unigame` framework are [`anyCards`](https://github.com/joshuaauerbachwatson/anyCards) and [`tictactoe`](https://github.com/joshuaauerbachwatson/tictactoe).
