<h1>oracleVpsNote</h1>
<hr>

###ubuntu
```shell
iptables -I INPUT 6 -m state --state NEW -p tcp --dport 80 -j ACCEPT
netfilter-persistent save
```
###listen status
```shell
sudo netstat -tlnp | grep nginx
```
