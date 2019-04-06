# Keyserver

This is a basic `keyserver` for cosmos-sdk applications. It contains the following routes:

```
GET     /version
GET     /keys
POST    /keys
GET     /keys/{name}?bech=acc
PUT     /keys/{name}
DELETE  /keys/{name}
POST    /tx/sign
```

`/tx/sign` is currently unimplemented.
