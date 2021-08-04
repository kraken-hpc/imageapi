# Security Considerations

The ImageAPI service must be run as root, and exposes powerful and dangerous methods through the API.

The current ImageAPI does not support direct authentication, though this may be implemented in the future.  As such, it's important to be careful with how the communication protocals are configured and protected.  

The current version of ImageAPI supports three protocols for communication:

- http
- https
- unix domain sockets

If ImageAPI needs to communicate only locally, we recommend using `unix` as the connection scheme.  This isolates the API from the network, and allows controlling access to the API by controlling permissions on the socket.

If ImageAPI needs network connectivity, the optimal scheme is to use HTTPS with carefully configured certificates and with the `--tls-ca=` option configured, which will enable mTLS authentication.  In this scheme, only clients with a valid, signed client certificate will be allowed to communicate with the API.

The http protocol should only be used in test environments.

See `imageapi-server -h`, or [running](running.md) for more details on runtime options.
