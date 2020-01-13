### Setup

#### Icinga 2 certificate

Generate cert on master
```
icinga2 pki new-cert  \
      --cn localhost   \
      --csr certs/localhost.csr    \
      --cert certs/localhost.crt   \
      --key certs/localhost.key
```

Sign on master
```
icinga2 pki sign-csr \
      --csr certs/localhost.csr \
      --cert certs/localhost.crt
```

#### Build switchboard & plugins

Switchboard
```
go build -o switchboard ./app
```

Connection plugin
```
go build -buildmode=plugin -o connection.plugin ./plugins/connection/
```

Check result plugin (optional)
```
go build -buildmode=plugin -o checkresult.plugin ./plugins/checkresult
```

PKI request plugin (optional)
```
go build -buildmode=plugin -o pki-request.plugin ./plugins/pki-request
```

#### Run the switchboard

```
CONNECTION_ADDR="localhost:5665" ./switchboard
```