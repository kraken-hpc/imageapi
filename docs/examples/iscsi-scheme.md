# Scheme: iSCSI

The following definition is an example of using iSCSI as a read-only attachment type, with `ext4` on the iSCSI object, and using Overlay to make the image read-write.  This has shown to be a highly scalable scheme for large-scale clusters.  Additionally, using `ext4` as a filesystem, while less effecient that `squashfs`, means that the image can be modified over time, which simplifies some workflows.

We recommend, when using this scheme, that the `mapped_lun` for the iSCSI target be specified with `write_protect=1`.

## Creating the image

As an example, we can use `targetcli` and the `fileio` backend to create the image. Generally, you would want to script this process.

1. Create the sparse image file:

   We'll use a sparse image so that we don't take a fixed amount of space:

   ```console
   dd if=/dev/zero of=/tmp/image.ext4 bs=1 count=0 seek=10G
   ```

   This will create the file `/tmp/image.ext4` with a maximum of `10G` of space, but initial usage of `0`.
2. Attach loopback file:

   ```bash
   $ losetup -f --show /tmp/image.ext4
   /dev/loop0
   ```

   This attached the file to `/dev/loop0`
3. Create the filesystem on the image:

   Now, create the filesystem:

   ```console
   mkfs.ext4 -Lroot /dev/loop0
   ```

4. Mount the filesystem:

   Create a mountpoint and mount the filesystem:

   ```console
   mkdir /mnt/image
   mount /dev/loop0 /mnt/image
   ```

5. Copy image/create image:

   At this point we can populate the filesystem with an image.  This could be already created, and we copy the files into place, or perhaps using something like `dnf --installroot=/mnt/image ...` to build the image.

6. Cleanup mount/attachment:

   ```console
   unmount /mnt/image
   losetup -d /dev/loop0
   ```

7. Create fileio backend store:

   ```console
   targetcli /backstores/fileio create file_ro_dev="/tmp/image.ext4" name="image.ext4"
   ```

8. Create iqn:

   ```console
   # targetcli /iscsi create iqn.2003-01.org.linux-iscsi.localhost.x8664
   iqn.2003-01.org.linux-iscsi.localhost.x8664:sn.XXXXXXXXXXXX
   ```

   Note the full IQN.

9. Create mapped_lun:

   ```console
   # targetcli /iscsi/iqn.2003-01.org.linux-iscsi.localhost.x8664:sn.XXXXXXXXXXXX/tpg1/luns create /backstores/fileio/image.ext4
   Created LUN 0.
   ```

   Note the LUN number in the output.

10. Create acl/lun mapping:

   ```console
   targetcli /iscsi/iqn.2003-01.org.linux-iscsi.localhost.x8664:sn.XXXXXXXXXXXX/tpg1/acls create iqn.2003-01.org.linux-iscsi.localhost:888
   targetcli /iscsi/iqn.2003-01.org.linux-iscsi.localhost.x8664:sn.XXXXXXXXXXXX/tpg1/acls/iqn.2003-01.org.linux-iscsi.localhost:888 create 0 0 write_protect=1
   ```

## Container definition

```json
{
  "name": "image.ext4",
  "command": "/usr/lib/systemd/systemd",
  "mount": {
    "kind": "overlay",
    "overlay": {
      "lower": [
        {
          "kind": "attach",
          "attach": {
            "fs_type": "ext4",
            "mount_options": [
              "ro"
            ],
            "attach": {
              "kind": "iscsi",
              "iscsi": {
                "lun": 0,
                "host": "192.168.3.253",
                "initiator": "iqn.2003-01.org.linux-iscsi.localhost:888",
                "target": "iqn.2003-01.org.linux-iscsi.localhost.x8664:sn.XXXXXXXXXXXX"
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
