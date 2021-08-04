# ImageAPI

The ImageAPI describes a restful (Swagger/OpenAPI 2.0) interface for attching, mounting, and launching system image containers.

## Overview

The ImageAPI is desgined to provide a flexible and efficient mechanism to deploy system images to stateless clusters.

This service is likely much more useful when combined with a tool like [Kraken](https://github.com/hpc/kraken) that can automate the image attach/load process in conjuction with network booting.

Unlike traditional network booting solutions, the ImageAPI is intended to be run as a persistent service that can load and unload system images as needed.  Images are run as privileged, namespaced processes, a kind of light-weight linux "container."

The API specification is contained in [swagger.yaml](swagger.yaml) .  Human-readable docs can be browsed at [docs/api/index.html](docs/api/index.html) .

It can also be browsed on [SwaggerHub](https://app.swaggerhub.com/apis/jlowellwofford/image-api/0.2.0).

## Architecture

### Object types

The ImageAPI understands three types of objects that work together:

- An **Attachment** is an object that describes a way of acquiring a block device file.  The product of an attachment is always a block device file (e.g., `/dev/rbd0`, `/dev/sda`, `/dev/loop0`, etc.). Current attachment types are:
  - **iscsi** attach an iSCSI network block device (experimental). See [docs/attachments/iscsi](docs/attachments/iscsi.md)
  - **local** attach an already-existing local block device (e.g. a local disk). See [docs/attachments/local](docs/attachments/local.md)
  - **loopback** attach a file in the root filesystem or on an existing mount as a loopback block device. See [docs/attachments/loopback](docs/attachments/loopback.md)
  - **rbd** attach a ceph/RBD network block device. See [docs/attachments/rbd](docs/attachments/rbd.md)
- A **Mount** describes any mount-type mechanism that provides a filesystem.  Current Mount types are:
  - **attach** mount a filesystem on a block device provided by one of the **Attachment** types. See [docs/mounts/attach](docs/mounts/attach.md)
  - **bind** bind-mount a directory contained in root (`/`) or an existing mount. See [docs/mounts/bind](docs/mounts/bind.md)
  - **nfs** mount a Network File System (NFS) mountpoint. See [docs/mounts/nfs](docs/mounts/nfs.md)
  - **overlay** create a read-write overlayfs mount over existing mountpoints (typically read-only). See [docs/mounts/overlay](docs/mounts/overlay.md)
- A **Container** describes a runnable system image on an existing mountpoint.  Containers need to be provided an init process and some optional occompanying meta-data.  SystemD is supported as an init processes as long as `systemd: true` is provided as an option to the container.  

### Object structure & nesting

In general, objects can depend on the existence of other objects.  For instance, Containers always need to reference a Mount object.  Mount objects may need to reference an Attach object or another Mount object.  This can be specified in one of two ways.  Objects another object depends on can either be referenced by ID, in which case that object depended on needs to already exist, or, object definitions can be nested.  In this way, a container definition along with the dependent mount and attach definitions can all be contained in one description.

If objects are provided in a nest fashion, dependencies will automatically be collected if the main object is removed, unless additional references are added.

### Some example object structures

Object structures are kept abstract to allow maximum flexiblity.  Objects can be nested in any way they can fit together, and objects can have an arbitrarily deep dependency chain.  Here are a few examples of valid, useful object structures:

- **Read-Write NFS attach** [Container -> NFS Mount (read-write)] Only really useful for a single image per node configuration.
- **Read-Only NFS with overlay** [Container -> Overlay Mount (read-write) -> NFS Mount (read-only)] This allows an extracted NFS shared image to be used by many nodes.  Tests show this is not the most efficient method, but it is easy to manage for small clusters.
- **Loopback filesytem on NFS mount** [Container -> Overlay Mount (read-write) -> Attach mount (read-only) -> Loopback file (e.g SquashFS, read-only) -> NFS Mount (read-only)].  This may seem like a lot of layers, but it leads to a way to distribute images much more efficiently than the previous to a larger cluster.
- **Read-write iSCSI/RBD attachment** [Container -> Attach Mount (e.g. ext4, read-write) -> iSCSI or RBD Attach (read-write)]. Only useful for single image per node configuration.  Like **Read-Write NFS attach**, but with block-level transfer instead of filesystem-level (more efficient in tests).
- **Read-only iSCSI/RBD attachment with overlay** [Container -> Overlay Mount (read-write) -> Attach Mount (e.g. SquashFS, read-only) -> RBD or iSCSI attach (read-only)] Provides a highly-scalable, multiple-node-per-image scheme for larger clusters.  Tests (publication pending) have indicated that RBD + SquashFS + Overlay provides the highest efficiency for scale in most configurations.

While this list provides some sensible combinations of objects, any object combination that can be described will work.  For instance, if there were a reason to do so, the following would be valid: [Container -> Overlay Mount -> Bind Mount -> Attach Mount -> Loopback Attach -> Attach Mount -> iSCSI Attach].  Of course, it's unlikely that this scheme would provide much efficiency.

### IDs and References

Every object is automatically provided a globally unique `ID` which can be used to operate on that object.  Additionally, Containers can be provided an (optional) unique `Name` which can be used in place of the `ID`.

All objects track internal references to that object.  When the reference count for an object reaches `0`, it will be garbage collected automatically.  Manually created mounts and attachments will automatically have a reference of `1`.  Mounts and attachments with non-zero reference counts can be forced removed (assuming they are not registered in the system as busy) with the `force=true` query option.

## Documentation

For an index of additional documentation, see: [docs/index](docs/index.md)