Deploy a complete demo by the cmd below
```shell
chmod +x statrup.sh

./startup.sh
```

**Note the default account of EMQX dashboard has been set to admin/admin**

The shell deploys EMQX5 enterprise by default, it also supports to deploy other versions by passing an additional arg `emqx4`, `emqx4-ee` and `emqx5`.

For example:
```shell
# deploy EMQX5 opensource version
./startup.sh emqx5
```