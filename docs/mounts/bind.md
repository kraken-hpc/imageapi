# Bind Mount

The Bind Mount object describes a bind mount into either the root (`/`) filesystem, or an existing Mount object.  Bind mounts provide a way to treat a subdirectory as an image.

The parameter `base: [ "root", "mount" ]` controls whether the bind mount should be located in the root filesystem or an existing mount.

Example specification of a bind mount (`base=root`):

```json
{
  "kind": "bind",
  "bind": {
          "path": "/chroots/image1",
          "base": "root",
          "recursive": true,
          "ro": true
  }
}
```
