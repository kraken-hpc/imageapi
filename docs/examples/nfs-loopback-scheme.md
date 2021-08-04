# Scheme: NFS + Loopback

While the [NFS](nfs-scheme.md) provides a lot of convenience because it uses NFS, which is very common, and simply serves images as subdirectories in an export, it tends to perform poorly for large-scale clusters.  We can maintain some of these convenience and significantly improve scalibility by storing loopback filesystems on NFS exports and loopback mounting those filesystems.  This maintains the use of NFS, but has performance characteristics more similar to those of the network block device attachment schemes.  NFS+Loopback is a reasonable choice for a mid-level cluster where there is a desire to avoid the more specialized network block device technologies, like iSCSI and Ceph/RBD.  We will illustrate using SquashFS as the filesystem.  A similar procedure to the one in [Scheme: iSCSI](iscsi-scheme.md) can be followed for other filesystem types.

## Creating the image

1. Create the NFS export

   We first need an NFS exported area.  Edit `/etc/exports` and add an entry like:

   ```file
   /images 192.168.3.0/24(ro)
   ```

   export with `exportfs -a` and verify the export with `showmount -e`.

   Note that, unlike the NFS scheme, we do not require `no_root_squash`.

2. Create image chroot

   Follow some procedure that provides a subdirectory containing the image files, e.g. a `dnf --installroot=... ...`, or perhaps a container build and mount.

3. Make the squashfs image

   From inside the image chroot directory, create a squashfs image with:

   ```console
   mksquashfs . /images/image.sqsh -noI -noX -noappend
   ```

   The image is now available to the cluster.

## Container definition

```json
{
  "name": "image.sqsh",
  "command": "/usr/lib/systemd/systemd",
  "state": "running",
  "systemd": true,
  "mount": {
    "kind": "overlay",
    "overlay": {
      "lower": [
        {
          "kind": "attach",
          "attach": {
            "fs_type": "squashfs",
            "mount_options": [
              "ro"
            ],
            "attach": {
              "kind": "loopback",
              "loopback": {
                "path": "/image.sqsh",
                "base": "mount",
                "mount": {
                  "kind": "nfs",
                  "nfs": {
                    "path": "/images",
                    "host": "192.168.3.253",
                    "ro": true
                  }
                }
              }
            }
          }
        }
      ]
    }
  }
}
```
