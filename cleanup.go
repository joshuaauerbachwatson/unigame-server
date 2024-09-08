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

// Handles periodic cleanup of players and games that have been idle long enough

package main

import (
	"fmt"
	"time"
)

// Cleanup function, expected to be invoked at regular intervals.
// The constant `playerTimeout` expresses how many times a player can be found
// unresponsive by this function before it is removed.   The player idle count is
// zeroed everytime the player app responds with a pong to a websocket ping from the server.
// So, the playerTimeout should be a small multiple of the pong wait time.
// The constant `gameFormationTimeout` expresses how many times a game may be found
// incomplete (number of players unknown or not as many players as needed).  Once a game
// is complete, its idle count is no longer used.  It will be deleted when it has no more
// players.
func cleanup() {
	// Note: deletion from a map in the scope of a 'range' loop is said to be safe:
	// https://stackoverflow.com/questions/23229975/is-it-safe-to-remove-selected-keys-from-map-within-a-range-loop
	for gameToken, game := range games {
		// Time out any games that have taken too long to find enough players
		if game.NumPlayers == 0 || len(game.Players) < game.NumPlayers {
			// Game not yet fully assembled, so subject to time limit
			game.IdleCount++
			if game.IdleCount > gameFormationTimeout {
				fmt.Printf("cleanup deleting incomplete game '%s' that has passed its time limit\n", gameToken)
				for _, player := range game.Players {
					if player.Client != nil {
						player.Client.Destroy()
					}
				}
				delete(games, gameToken)
				continue
			}
		} else {
			game.IdleCount = 0
		}
		// Timeout any players that have been idle too long.  Delete the game if it has no player.
		for playerOrder, player := range game.Players {
			player.IdleCount++
			if player.IdleCount > playerTimeout {
				fmt.Printf("cleanup deleting player %d from game %s\n", playerOrder, gameToken)
				if player.Client != nil {
					player.Client.Destroy()
				}
				delete(game.Players, playerOrder)
			}
		}
		if len(game.Players) == 0 {
			fmt.Printf("cleanup discarding game %s because it no longer has any players\n", gameToken)
			delete(games, gameToken)
		}
	}
}

// Start a ticker to do cleanup every 'cleanupPeriod' seconds
func startCleanupTicker() {
	ticker := time.NewTicker(cleanupPeriod * time.Second)
	go func() {
		for {
			<-ticker.C
			cleanupCounter++
			cleanup()
		}
	}()
}
