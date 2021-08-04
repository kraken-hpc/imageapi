# Scheme: RBD

The RBD scheme is similar to the [Scheme: iSCSI](iscsi-scheme.md) in structure, as both are network block device attachment schemes.  In our testing, RBD + SquashFS + Overlay provided the most scalable results.   Moreover, the added features of Ceph/RBD like replication and failover, make this a promising solution for large-scale, demanding clusters.

Note: in this example we use the admin account to attach the RBD for simplicity.  You almost certainly want a more restricted account for this use.

## Creating the image

We will demostrate using SquashFS here.  Other filesystem types would follow by running `mkfs.*` against the rbd device, then copying files over (or writing them in place).

1. Create image chroot

   Follow some procedure that provides a subdirectory containing the image files, e.g. a `dnf --installroot=... ...`, or perhaps a container build and mount.

2. Create the RBD object

   ```console
   rbd create --image-feature layering --image-shared -s 100G image1
   ```

3. Map the RBD object to the image creation machine

   ```console
   # rbd map --monitor 192.168.3.253 --pool rbd --image iamge1 --id admin --secret XXXXXXXXXXXX
   /dev/rbd0
   ```

4. Write the image

   From within the chroot directory, build the squashfs directly to the rbd device:

   ```console
   mksquashfs . /dev/rbd0 -noI -noX -noappend
   ```

5. Cleanup attachment

   ```console
   rbd unmap -d 0
   ```

   The image is now ready to use.

## Container definition

```json
{
  "name": "test",
  "command": "/usr/lib/systemd/systemd",
  "mount": {
    "kind": "overlay",
    "overlay": {
      "lower": [
        {
          "kind": "attach",
          "attach": {
            "kind": "rbd",
            "fs_type": "squashfs",
            "mount_options": [
              "ro"
            ],
            "attach": {
              "kind": "rbd",
              "rbd": {
                "image": "image1",
                "monitors": [
                  "192.168.3.253"
                ],
                "options": {
                  "name": "admin",
                  "ro": true,
                  "secret": "XXXXXXXXXXXX"
                },
                "pool": "rbd"
              }
            }
          }
        }
      ]
    }
  },
  "state": "running",
  "systemd": true
}
```
