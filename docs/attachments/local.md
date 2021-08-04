# Local Attach

The local attachment type can be used to make a locally available block device that already exists available as an Attachment object for the ImageAPI.  This can allow, e.g. a local disk to be mounted.  This attachment type will fail if the device file does not exists, or if the device is not a block device.

Example of a Local attachment specification:

```json
{
  "kind": "local",
  "local": {
    "path": "/dev/sda"
  }
}
```
