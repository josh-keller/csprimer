# Proxy Server

## Goal

Write a proxy that can send a request on to an upstream server and return the server's response.

## Setup

1. Need static web server
2. Run proxy
3. Make curl requests
4. Make browser request

## Approach

1. Listen for connections
2. For each connection:
   1. Establish the connection
   2. Read raw byets from it
   3. Parse the headers (eventually use this to proxy to multiple backends by reading the host header)
      * host/path
      * body size
   4. After headers are parsed, open connection to upstream server
   5. Send message to upstream (keep same packets?)
   6. Receive from upstream
   7. Forward these directly back?
   8. Close connection (for now)


