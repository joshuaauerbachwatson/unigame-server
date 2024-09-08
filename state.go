/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Stores the (volatile, in memory) state of all the active games.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"golang.org/x/exp/maps"
)

// The state of one game
type Game struct {
	Players   map[uint32]*Player `json:"players"`   // key is the player's "order" number
	IdleCount int                `json:"idleCount"` // global idle count for the game as a whole
	// The game idle count is only used while the player list is incomplete.
	NumPlayers int `json:"numPlayers"` // The expected number of players for this game
	// Note: the number of players in the Players map should not exceed NumPlayers but may be less as
	// players join the game.  A NumPlayers value of 0 means "unknown", which may be case transiently.
	Hub *Hub // The Websocket "Hub" for the game (not serialized)
}

// The state of one Player
type Player struct {
	Token     string  `json:"token"`     // Player's token (encodes name and order number)
	IdleCount int     `json:"idleCount"` // Idle count for this player.
	Client    *Client // The Websocket "client" for the player (not serialized)
}

type DumpedState struct {
	CleanupCounter int              `json:"cleanupCounter"`
	Games          map[string]*Game `json:"games"`
}

// Map from game tokens to Game structures
var games = make(map[string]*Game)

// Counter for the number of times cleanup has run
var cleanupCounter int

// Subroutine to make the player list of a game
func makePlayerList(game *Game) string {
	keys := maps.Keys(game.Players)
	slices.Sort(keys)
	list := ""
	delim := ""
	for _, key := range keys {
		list += (delim + game.Players[key].Token)
		delim = " "
	}
	return strconv.Itoa(game.NumPlayers) + " " + list
}

// Handler for an admin function to dump the entire state of the server.
// This is an aid during development.  We might need something more sophisticated
// for observability in the long run.
func dump(w http.ResponseWriter) {
	ans := DumpedState{CleanupCounter: cleanupCounter, Games: games}
	encoded, err := json.MarshalIndent(ans, "", "  ")
	if err != nil {
		indicateError(http.StatusInternalServerError, err.Error(), w)
		return
	}
	fmt.Println("dump called")
	fmt.Println(string(encoded))
	w.Write(append(encoded, byte('\n')))
}

// Handler for an admin function to reset to the empty state
func reset() {
	fmt.Println("reset called")
	games = make(map[string]*Game)
	cleanupCounter = 0
}
