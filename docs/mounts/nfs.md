# NFS Mount

The NFS Mount object describes an NFS mount point.  The mount can be read-only if `ro: true` is specified.  Additionally, the NFS version can be controlled with the `version: <vers>` option.  Varous NFS tuning options can be provided thropugh `options: []`.

Example specification of an NFS mount:

```json
{
  "kind": "nfs",
  "nfs": {
          "path": "/var/chroots",
          "host": "192.168.3.253",
          "ro": true
  }
}
```

For a complete example using NFS to host images as subdirectories, see [../examples/nfs-scheme](../examples/nfs-scheme.md) .

For a complete example using NFS to host squashfs image files, see [../examples/nfs-loopback-scheme](../examples/nfs-loopback-scheme.md)
