# Overlay Mount

The Overlay Mount object provides a way, using OverlayFS, to present a read-only mount as a read-write mount.  This is useful in several differnt image schemes.

The Overlay Mount object supports layered mounts.  This is handled by providing multiple mount descriptors to the `lower: []` section.

Overlay Mounts are always read-write.

Example Overlay Mount specification:

```json
{
    "kind": "overlay",
    "overlay": {
      "lower": [
        {
            // Mount specification
            // ...
        }
      ]
    }
}
```

The following examples give complete usage cases for Overlay Mounts:

- [../examples/nfs-scheme](../examples/nfs-scheme.md)
- [../examples/nfs-loopback-scheme](../examples/nfs-loopback-scheme.md)
- [../examples/rbd-scheme](../examples/rbd-scheme.md)
- [../examples/iscsi-scheme](../examples/iscsi-scheme.md)
