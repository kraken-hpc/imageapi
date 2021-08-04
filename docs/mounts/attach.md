# Attach Mount

The Attach Mount mounts any attachment type.  In addition to the attachment specification, it contains a filesystem type and optional `mount_options`.

Example Attach Mount specification:

```json
{
    "kind": "attach",
    "attach": {
        "kind": "rbd",
        "fs_type": "squashfs",
        "mount_options": [
            "ro"
        ],
        "attach": { 
            // attachment specification
            // ...
        }
    }
}
```

See [../examples/iscsi-scheme](../examples/iscsi-scheme.md) or [../examples/rbd-scheme](../examples/rbd-scheme.md) for complete examples.