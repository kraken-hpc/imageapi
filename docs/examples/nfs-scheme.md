# Scheme: NFS

This is the simplest scheme we provide an example for. In this scheme, we NFS export a directory containing subdirectories that are image chroot environments.  We directly mount these directories (read-only) and overlay them.  While this is the simplest scheme, it is generally the slowest and least scalible.

## Creating the image

1. Create the NFS export

   We first need an NFS exported area.  Edit `/etc/exports` and add an entry like:

   ```file
   /images 192.168.3.0/24(ro,no_root_squash)
   ```

   export with `exportfs -a` and verify the export with `showmount -e`.

   Note that `no_root_squash` is require for this scheme to function.

2. Create image chroot

   Follow some procedure that provides a subdirectory containing the image files, e.g. a `dnf --installroot=... ...`, or perhaps a container build and mount.  Make sure this chroot is placed under the NFS export directory, e.g. `/images/image1`.

   The image is now available to the cluster.

## Container definition

```json
{
  "name": "image1",
  "command": "/usr/lib/systemd/systemd",
  "state": "running",
  "systemd": true,
  "mount": {
    "kind": "overlay",
    "overlay": {
      "lower": [
        {
          "kind": "nfs",
          "nfs": {
            "path": "/images/image1",
            "host": "192.168.3.253",
            "ro": true
          }
        }
      ]
    }
  }
}
```
