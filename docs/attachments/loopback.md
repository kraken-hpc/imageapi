# Loopback Attach

The Loopback Attach allows a file that lives on either the root filesystem (`/`) or an existing mount to be attached as a lopback block device.  This allows, e.g. a SquashFS image that exists on another mount to be attached and mounted.

The parameter, `base: [ "mount", "root" ]` controls if the file is to be found in the root filesystem or in a mount.

Example of a Loopback attachment specification (in the root filesystem):

```json
{
    "kind": "loopback",
    "loopback": {
    "path": "/rhel8.1-x86_64.sqsh",
    "base": "root",
}
```

To see an example use if loopback with an NFS mount, see [../examples/nfs-loopback-scheme](../examples/nfs-loopback-scheme.md)
