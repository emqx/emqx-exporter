Deploy a complete demo by the cmd below
```shell
chmod +x statrup.sh

./startup.sh
```

**Note the default account of EMQX dashboard has been set to admin/admin**

The shell deploys EMQX 5 enterprise by default, it also supports deploying other versions by passing an additional arg `emqx4`, `emqx4-ee` and `emqx5`.

For example:
```shell
# deploy EMQX 5 opensource version
./startup.sh emqx5
```