# ImageAPI

The ImageAPI describes a restful (Swagger/OpenAPI 2.0) interface for attching, mounting, and launching system image containers.

This is desgined to provide a flexible and efficient mechanism to deploy system images to stateless clusters.

This service is likely much more useful when combined with a tool like [Kraken](https://github.com/hpc/kraken) that can automate the image attach/load process in conjuction with network booting.

The API specification is contained in [swagger.yaml](swagger.yaml) .

It can also be browsed on [SwaggerHub](https://swaggerhub.com).

# Example interaction

1. Make sure that the `rbd` and `overlay` modules are loaded:
   ```bash
   modprobe rbd overlay
   ```
2. Start the `imageapi-server` service.  We'll just run it by hand:
   ```bash
   $ nohup sudo ./imageapi-server --port 8080 --scheme http &
   ```
   This starts and backgrounds the service on `127.0.0.1:8080` .  This is insecure, but good for testing.
3. Attach an RBD object.  We will attach one named `systemd.sqsh` that already contains a `systemd` based image in a `squashfs` filesystem.

   ```bash
   $ curl -s -XPOST -H 'Content-Type: application/json' -d '{"monitors":["192.168.1.48"],"pool":"rbd","image":"systemd.sqsh","options":{"ro":true,"name":"admin","secret":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}}' http://localhost:8080/attach/rbd | jq
   {
     "image": "systemd.sqsh",
     "monitors": [
       "192.168.1.48"
     ],
     "options": {
       "name": "admin",
       "ro": true,
       "secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
     },
     "pool": "rbd"
   }
   ```
   We have successfully attached the RBD object (read-only).  We can see this with:
   ```bash
   $ dmesg | tail -n3
   [ 9584.859653] libceph: mon0 (1)192.168.1.48:6789 session established
   [ 9584.861001] libceph: client124116 fsid 4bc8e219-dd58-473f-ac79-c947bd1c4d41
   [ 9584.905617] rbd: rbd0: capacity 10737418240 features 0x1
   $ ls -l /dev/rbd0
   brw-rw---- 1 root disk 253, 0 Feb 10 23:53 /dev/rbd0
   ```
4. Now that we have attached the object, we need to mount it.
   ```bash
   $ curl -s -XPOST -H 'Content-type: application/json' -d '{"id": 0, "fs_type": "squashfs", "mount_options": [ "ro" ] }' http://localhost:8080/mount/rbd | jq
   {
     "fs_type": "squashfs",
     "id": 0,
     "mount_options": [
       "ro"
     ],
     "mountpoint": "/var/run/imageapi/mounts/mount_746184320"
   }
   ```
   We see that the image got mounted under `/var/run/imageapi/mounts/mount_746184320`.
   ```bash
   $ sudo ls /var/run/imageapi/mounts/mount_746184320
   bin  boot  dev	etc  home  lib	lib64  media  mnt  opt	proc  root  run  sbin  srv  sys  tmp  usr  var
   ```
5. To make make a non-destructive, locally read-write image, we mount an `overlay` over this image.
   ```bash
   $ curl -s -XPOST -H 'Content-type: application/json' -d '{ "lower": [ 0 ]}' http://localhost:8080/mount/overlay | jq
   {
     "id": 1,
     "lower": [
       0
     ],
     "mountpoint": "/var/run/imageapi/mounts/mount_725692383",
     "upperdir": "/var/run/imageapi/mounts/upper_014522290",
     "workdir": "/var/run/imageapi/mounts/work_400792425"
   }
   ```
   Note that the `"lower"` array referenced the rbd mount `"id"` from above.  We can handle multiple overlay layers by adding more entries to this list.

   The `mountpoint` at `/var/run/imageapi/mounts/mount_725692383` is read-write.
6. Now we can define our container on this mount:
   ```bash
   $ curl -s -XPOST -H 'Content-type: application/json' -d '{ "mount": { "id": 1, "kind": "overlay" }, "command": "/usr/lib/systemd/systemd", "state": "created", "systemd": true }' http://localhost:8080/container | jq
   {
     "command": "/usr/lib/systemd/systemd",
     "logfile": "/var/run/imageapi/logs/0-1613002114.log",
     "mount": {
       "id": 1,
       "kind": "overlay"
     },
     "namespaces": null,
     "state": "created",
     "systemd": true
   }
   ```
   The `mount` structure specified using an `overlay` mount of `id: 1` that we just created.  `command` specifies the entrypoint command, in this case, `systemd`.  If we're running `systemd` some extra container setup is necessary.  The `"systemd": true` option makes sure that happens.

   Note, `namespaces` is currently unused.  All containers get `mount`, `pid`, `uts`, and `ipc` namespaces by default.

   We see a reference to a log file for this container.  Currently it's pretty boring:
   ```bash
   $ sudo cat /var/run/imageapi/logs/0-1613002114.log
   2021/02/11 00:08:34 container(0): container created
   ```
   That's because we didn't request that the container actually start.  We could have with `"state": "running"`.
7. Since we didn't auto-start our container, let's start it:
   ```bash
   $ curl -s -XGET http://localhost:8080/container/0/running | jq
   {
     "command": "/usr/lib/systemd/systemd",
     "logfile": "/var/run/imageapi/logs/0-1613002114.log",
     "mount": {
       "id": 1,
       "kind": "overlay"
     },
     "namespaces": null,
     "state": "running",
     "systemd": true
   }
   ```

   We can `ps` to see the processes running.

   ```bash
   $ ps -elF --forest
   ...
   4 S root        1074    1010  0  80   0 - 88261 -       7756   0 21:36 pts/0    00:00:00  |           \_ sudo ./imageapi-server --scheme=http --port=8080
   4 S root        1075    1074  0  80   0 - 178813 -     23640   0 21:36 pts/0    00:00:00  |               \_ ./imageapi-server --scheme=http --port=8080
   4 S root        1534    1075  2  80   0 - 24282 -      13004   0 21:44 ?        00:00:00  |                   \_ /usr/lib/systemd/systemd
   4 S root        1563    1534  0  80   0 -  7761 -      10572   0 21:44 ?        00:00:00  |                       \_ /usr/lib/systemd/systemd-journald
   4 S root        1610    1534  0  80   0 -  5176 -       9408   1 21:44 ?        00:00:00  |                       \_ /usr/lib/systemd/systemd-logind
   ...
   ```

   We can use `nsenter` to "enter" the container:
   ```bash
   $ sudo nsenter -t 1534 -a bash
   [root@kraken /]# ps waux
   USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
   root           1  0.5  0.4 170732 12604 ?        Ss   00:14   0:00 /usr/lib/systemd/systemd
   root          30  0.0  0.5  41236 12832 ?        Ss   00:14   0:00 /usr/lib/systemd/systemd-journald
   root          73  0.0  0.4 248184 11180 ?        Ss   00:14   0:00 /usr/sbin/sssd -i --logger=files
   root         103  0.0  0.5 251732 12928 ?        S    00:14   0:00 /usr/libexec/sssd/sssd_be --domain implicit_files --uid 0 --gid 0 --logger=files
   root         109  0.0  1.5 274916 38160 ?        S    00:14   0:00 /usr/libexec/sssd/sssd_nss --uid 0 --gid 0 --logger=files
   root         170  0.2  0.3  20704  9576 ?        Ss   00:16   0:00 /usr/lib/systemd/systemd-logind
   root         194  0.3  0.1 231648  4200 ?        S    00:16   0:00 bash
   root         206  0.0  0.1 234240  3760 ?        R+   00:16   0:00 ps waux
   ```
   We can check our currently running containers:
   ```bash
   $ curl -s http://localhost:8080/container | jq
   [
     {
       "command": "/usr/lib/systemd/systemd",
       "logfile": "/var/run/imageapi/logs/0-1613002114.log",
       "mount": {
         "id": 1,
         "kind": "overlay"
       },
       "namespaces": null,
       "state": "running",
       "systemd": true
     }
   ]
   ```

8. Finally, let's tear it all down:
   ```bash
   $ curl -XDELETE http://localhost:8080/container/0
   $ curl -XDELETE http://localhost:8080/mount/overlay/1
   $ curl -XDELETE http://localhost:8080/mount/rbd/0
   $ curl -XDELETE http://localhost:8080/attach/rbd/0
   $ sudo cat /var/run/imageapi/logs/0-1613002114.log
   2021/02/11 00:08:34 container(0): container created
   2021/02/11 00:14:54 container(0): starting container
   2021/02/11 00:14:54 container(0): validating image
   2021/02/11 00:14:54 container(0): validating init
   2021/02/11 00:14:54 init: making all mounts private
   2021/02/11 00:14:54 init: preparing image
   2021/02/11 00:14:54 init: mounting /proc
   2021/02/11 00:14:54 init: mounting /dev
   2021/02/11 00:14:54 init: mounting /dev/shm
   2021/02/11 00:14:54 init: mounting /dev/mqueue
   2021/02/11 00:14:54 init: mounting /dev/pts
   2021/02/11 00:14:54 init: mounting /sys
   2021/02/11 00:14:54 init: mounting /run
   2021/02/11 00:14:54 init: mounting /tmp
   2021/02/11 00:14:54 init: mounting /sys/fs/cgroup
   2021/02/11 00:14:54 init: mounting /var/lib/journal
   2021/02/11 00:14:54 init: making device file /dev/null
   2021/02/11 00:14:54 init: making device file /dev/zero
   2021/02/11 00:14:54 init: making device file /dev/full
   2021/02/11 00:14:54 init: making device file /dev/tty
   2021/02/11 00:14:54 init: making device file /dev/random
   2021/02/11 00:14:54 init: making device file /dev/urandom
   2021/02/11 00:14:54 init: creating symlink /dev/pts/ptmx -> /dev/ptmx
   2021/02/11 00:14:54 init: creating symlink /proc/self/fd -> /dev/fd
   2021/02/11 00:14:54 init: creating symlink /proc/self/fd/0 -> /dev/stdin
   2021/02/11 00:14:54 init: creating symlink /proc/self/fd/1 -> /dev/stdout
   2021/02/11 00:14:54 init: creating symlink /proc/self/fd/2 -> /dev/stderr
   2021/02/11 00:14:54 init: executing init
   2021/02/11 00:20:07 container(0): container deleted
   $ ls -l /dev/rbd0
   ls: cannot access '/dev/rbd0': No such file or directory

   ```
