# RBD Attach

The RBD Attach allows a Ceph/RBD object to be attached.  This has proven to be one of the most efficient attachment mechanisms, and is one of the better tested mechanisms.

The RBD Attach only supports RBD protocol v1, since that is supported at the kernel level.  In some cases, this protocal may need to be enabled in your Ceph deployment.

The RBD Attach supports the specification of multiple monitors (for failover and load balancing), rbd namespaces, snapshots, and a number of different RBD tuning parameters.

If many nodes are going to use the same (read-only) RBD object, that object must be created with the `--image-shared` option.

Tests have found that the most efficient and supported image attach scheme is SquashFS over RBD with an Overlay mount.

***Important: specifying an RBD attach requires specifying an RBD secret.  Be careful which secrets this could potentially expose.***

Example RBD attachment specification:

```json
{
    "kind": "rbd",
    "rbd": {
        "image": "centos-slurm",
        "pool": "rbd",
        "monitors": [
            "192.168.3.253"
        ],
        "options": {
            "name": "user",
            "ro": true,
            "secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
        }
    }
}
```

For an example of using RBD + SquashFS + Overlay, see [../examples/rbd-scheme](../examples/rbd-scheme.md).

This scheme is also used in the Kraken/Layercake "layerk8s" example.  See [Kraken/layercake](https://github.com/kraken-hpc/kraken-layercake).