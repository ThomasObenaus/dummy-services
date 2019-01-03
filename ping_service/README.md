# Ping-Service

The ping-service is a simple service for testing purposes.
When you send a request to it's endpoint, the service tries to forward this request to other instances of the ping-service. This is done for a defined number of hops or "pings". For each hop a "ping" is added to the response. The last receiver in the chain stops forwarding and adds a "pong" to concatenated message list.
