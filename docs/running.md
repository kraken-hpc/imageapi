# Running the ImageAPI service

The ImageAPI service has the following runtime options;

```console
imageapi-server -h
Usage:
  imageapi-server [OPTIONS]

This API specification describes a service for attaching, mounting and preparing container images and manipulating those containers.

In general, higher level objects can either reference lower level objects (e.g. a mount referencing an attachment point) by a reference ID,
or, they can contain the full specification of those lower objects.

If an object references another by ID, deletion of that object does not effect the underlying object.

If an object defines a lower level object, that lower level object will automatically be deleted on deletion of the higher level object.

For instance, if a container contains all of the defintions for all mount points and attachments, deletion of the container will automatically unmount
and detach those lower objects.


Application Options:
      --scheme=            the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=   grace period for which to wait before killing idle connections (default: 10s)
      --graceful-timeout=  grace period for which to wait before shutting down the server (default: 15s)
      --max-header-size=   controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body. (default: 1MiB)
      --socket-path=       the unix socket to listen on (default: /var/run/imageapi.sock)
      --host=              the IP to listen on (default: localhost) [$HOST]
      --port=              the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=      limit the number of outstanding requests
      --keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download) (default: 3m)
      --read-timeout=      maximum duration before timing out read of the request (default: 30s)
      --write-timeout=     maximum duration before timing out write of the response (default: 60s)
      --tls-host=          the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=          the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=   the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=           the private key to use for secure connections [$TLS_PRIVATE_KEY]
      --tls-ca=            the certificate authority file to be used with mutual tls auth [$TLS_CA_CERTIFICATE]
      --tls-listen-limit=  limit the number of outstanding requests
      --tls-keep-alive=    sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)
      --tls-read-timeout=  maximum duration before timing out read of the request
      --tls-write-timeout= maximum duration before timing out write of the response

Help Options:
  -h, --help               Show this help message
```

In general, you will want to specify at least `--scheme`.  If using, `--scheme=https` you will want to specify the relevant `--tls-*` options.  You will likely also want to specify the `--host` and `--port` options.

To see some security-based recommendations on how to run ImageAPI, see [security](security.md).
