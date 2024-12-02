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

package main

// Constants used in the backend

const (
	// Minimum Length of a "game token" in characters.  Note that what we call a game token here is a composite
	// formed by the appid and what the client side calls a game token, separated by an underscore.  The
	// minimum length here assumes that both the app id and the game token are minimum length.
	minGameTokenLen = 12

	// URL paths representing verbs
	pathReset     = "/reset"
	pathDump      = "/dump"
	pathWebsocket = "/websocket"

	// Period at which the cleanup function is called, in seconds
	cleanupPeriod = 15

	// Default port to listen on if a port is not specified via the environment
	defaultPort = "80"

	// Timeouts for idle players and for game formation, as multiples of the cleanup period.
	playerTimeout        = 90 / cleanupPeriod
	gameFormationTimeout = 300 / cleanupPeriod

	// Query value keys used for websocket creation
	playerKey     = "Player"
	gameTokenKey  = "GameToken"
	numPlayersKey = "NumPlayers"
)
