This server can handle map operation requests.
The server consumes the messages from SQS, parses them, operates
them on in memory concurrent-map and rights them to log file.
To run the server (and client) it's necessary to change the url 
of the SQS to an existing SQS url.

My assumptions - 
SQS request takes much more time than execute map operation
(and all what that includes).
Messages order and duplicate messages are less critical than 
performance.
