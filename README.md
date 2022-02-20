This server can handle map operation requests.
The server consumes the messages from SQS, parses them, operates
them on in memory map (sorted by insert order) and rights them to log file.
To run the server (and client) it's necessary to change the url 
of the SQS to an existing SQS url.
