# iSCSI Attach

***iSCSI attachments are currently considered experimental***

The iSCSI attach driver allows ImageAPI to attach an iSCSI lun.  Currently, no authentication mechanisms are supported. In addition to the basic parameters required to specify the target and initiator, a handful of iSCSI tuning parameters are also supported.  See the API documentation for details.

Example of a minimal iSCSI attachment specifcation:

```json
{
    "kind": "iscsi",
    "iscsi": {
        "lun": 0,
        "host": "192.168.3.253",
        "initiator": "iqn.2003-01.org.linux-iscsi.my-host:888",
        "target": "iqn.2003-01.org.linux-iscsi.my-host.x8664:sn.XXXXXXXXXXXX"
        }
}
```

To see an examle use of iSCSI with an attachment mount and an overlay, see [../examples/iscsi-scheme](../examples/iscsi-scheme.md)
