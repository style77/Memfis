# Memfis

Wannabe Redis database, that saves all the data only in it's memory. Fully encrypted data, everything is removed after closing the server.


## Setup

To get executable simply run `go build -o ./builds`


## Endpoints

#### Request
`POST /authorize`
```
{
    "username": "root",
    "password": "root"
    }
```

#### Response

```
{
    "token": "eDF5UrXadWYDvVtwTtN9"
}
```

#### Description
`username` and `password` should be set in config.env. This token is only used to obtain API key

---------

#### Request
`POST /token`
```
{
    "token":"eDF5UrXadWYDvVtwTtN9"
}
```

#### Response
```
{
    "expiring": 1659762047,
    "token": "BoBHapHb.1659762047.c3R5bGU=",
    "user": "style"
}
```

#### Description

`Expiring` is timestamp, when api key is going to stop working, the time depends on config `API_TOKEN_EXPIRE_TIME` option.
`token` is API Key that is needed to use to send requests to API
`user` is user that generated that key in `/authorize` endpoint

---------

#### Request
`GET /data`
```
{
    "token": "BoBHapHb.1659762047.c3R5bGU=", 
    "name": "var"
}
```

#### Response
```
{
    "formatted": "string;var1;pol",
    "name": "var",
    "type": "int",
    "value": 123
}
```

#### Description

`formatted` is just every field formatted to one string
`name` variable name that is saved in `Memfis`
`type` type of variable, right now it's detected by `Memfis` itself. \
It needs tests, because sometimes it's detecting wrong type.
`ALL TYPES: string, int, float, array, json`
`value` value of variable

---------

#### Request
`POST /data`
```
{
    "token": "BoBHapHb.1659762047.c3R5bGU=",
    "name": "var", 
    "value": 123
}
```

#### Response
```
{
    "formatted": "string;var1;pol",
    "name": "var",
    "type": "int",
    "value": 123
}
```

#### Description
Same description as in `GET /data`

---------

#### Request
`DELETE /data`
```
{
    "token": "BoBHapHb.1659762047.c3R5bGU=",
    "name": "var"
}
```

#### Response
```
{
    "formatted": ";;<nil>",
    "name": "",
    "type": "",
    "updated": 1659741404,
    "value": null
}
```

#### Description
`formatted` formatted string is a little bugged but at least shows what should show
`updated` timestamp of when data were deleted
everything else is same as in `GET /data`

---------

#### Request
`PATCH /data`
```
{
    "token": "BoBHapHb.1659762047.c3R5bGU=",
    "name": "var",
    "value": "123"
}
```

#### Response
```
{
    "formatted": "int;var;123",
    "name": "var",
    "type": "int",
    "updated": 1659741641,
    "value": 123
}
```

#### Description
`updated` timestamp of when data were updated
everything else is same as in `GET /data`
data's type is chaning automatically

---------