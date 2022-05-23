#### HTTP 加密传输

##### `RSA` 加密算法

`RSA`使用公钥对数据进行加密，使用私钥进行解密

`RSA`也可以使用私钥对数据进行签名，然后公钥验签

`RSA`加密算法已经被破解到了`768`位了，所以，就目前来说`RSA 1024`位应该不算是很安全的，所以一般都需要使用 `RSA 2048`位的

###### `AES`加密算法

`AES`有五种加密算法，最常用的就是`AES CBC`（密码分组连接模式）

##### `HTTP` 加密传输方法

`Request`：`Base64（ Aes-CBC-128 Encrypt(data)+RSA-1024 Encrypt(Aes-CBC-128-salt)+len(RSA-1024 Encrypt(Aes-CBC-128-salt) ）`

分步骤：

```
d1 := Aes-CBC-128 Encrypt(data)
d2 := RSA-1024 Encrypt(Aes-CBC-128 salt)
d3 := len(d2)
d4 := base64(d1+d2+d3)
```

`Response`：`Base64 (Aes-CBC-128 Encrypt(data)+RSA-1024 Sign(Aes-CBC-128-salt)+len(RSA-1024 Sign(Aes-CBC-128-salt))`

分步骤：

```
d1 := Aes-CBC-128 Encrypt(data)
d2 := RSA-1024 Sign(d1)
d3 := len(d2)
d4 := Base64(d1+d2+d3)
```

`License API`设计

通过`pn`和`licenseId`获取`License`

```
/licenses/pns/:pn/licenseIds/:licenseId
```

通过`serviceName`和`licenseId`获取`License`

```
/licenses/service-names/:serviceName/licenseIds/:licenseId
```

