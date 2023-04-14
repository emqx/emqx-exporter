# Systemd Unit

If you are using distribution packages or the copied repository, you don't need to deal with these files!

The unit files (`*.service` and `*.socket`) in this directory are to be put into `/etc/systemd/system`.
It needs a user named `emqx_exporter`, whose shell should be `/sbin/nologin` and should not have any special privileges.
It needs a sysconfig file in `/etc/sysconfig/emqx_exporter`.
A sample file can be found in `sysconfig.emqx_exporter`.
