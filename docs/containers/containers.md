# Containers

Containers are the primary object of most ImageAPI definitions.  A container defines a runable mount.

In addition to a Mount specification to use, containers must at least define a `command:`.  This is the `init` command for the image.  In general, this should be some long-running initialization process.

Containers support using `systemd` as the init process, but in this case, the property `systemd: true` **must** be set.  The `systemd: true` property tells ImageAPI to do some extra initialization that SystemD requires, as well as use the appropriate signal to tell SystemD to shutdown.

Containers run with the full privileges of the `imageapi-service` (typically `root`).  However, they run inside a set of Linux namespaces.  While the container API allows for specification of which namespaces to use, these are currently ignored and the following namespaces are used by default: `Mount`, `PID`, `UTS` and `IPC`.  Containers are always run with the root of the mount appearing to be the root of the filesystem (via `Mount` namespaces).

Optionally, containers can be provided a unique `name:` property.  This allows controlling container lifecycle without tracking internal `id:` properties.

The states of either `created` or `running` can be provided as part of the container definition.  This will control whether ImageAPI will attempt to run the container on definition.

Once the container is started, the `logfile:` property will give a local path to a log file for that container.  Note console logs being printed to this file relies on logs being printed to standard sockets, which does not hold true, e.g. for systemd.

Example container definition:

```json
{
    "name": "simple-container",
    "command": "/init.sh",
    "state": "created",
    "mount": {
        // Mount specification
        // ...
    }
}
```

Exampel container definition (systemd):

```json
{
    "name": "systemd-container",
    "command": "/usr/lib/systemd/systemd",
    "systemd": true,
    "state": "running",
    "mount": {
        // Mount specification
        // ...
    }
}
```

See all of the examples container definitions complete with mount specifications.
