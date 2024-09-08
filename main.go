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

// The backend to the anyCards game, to be gradually evolved into something better.

package main

import (
	"fmt"

	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Main entry point
func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading the .env file: %v", err)
		return
	}

	// Set up handlers
	// Websocket initiation.  This should carry all traffic from the app itself
	http.Handle(pathWebsocket, EnsureValidToken()(
		http.HandlerFunc(newWebSocket),
	))

	// The Dump feature (requires admin role)
	http.Handle(pathDump, EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if body := screenRequest(w, r); body != nil {
				if isAdmin(w, r) {
					dump(w)
				}
			}
		}),
	))

	// The Reset feature (requires admin role)
	http.Handle(pathReset, EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if body := screenRequest(w, r); body != nil {
				if isAdmin(w, r) {
					reset()
				}
			}
		}),
	))

	// Permit port override (default 80)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Bind to port address
	bindAddr := fmt.Sprintf(":%s", port)
	fmt.Printf("==> Server listening at %s\n", bindAddr)

	// Start cleanup ticker
	startCleanupTicker()

	// Start serving requests
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	// No reasonable recovery at this point, just exit
	fmt.Println(err)
}
