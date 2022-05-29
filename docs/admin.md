# User

## Create User

**Request**
```http
POST /admin/user
Content-type: application/json
{
    "username": "fernando_kim",
    "password": "qiewgiqweijciqwe",
    "email": "qgiwjen.cwof",
    "accessControl": "SYSTEM",
}
```

- Data structure has been changed. Refer to [this](/docs/data-structure.md).

**Response**
```http
200 OK
Content-type: application/json
{
    "status": true,
    "data": {"id": "12040912512703foq", "count": 1}
}
```


## Read User
**Request** 

```http
GET /admin/user?sort=username_asc&limit=10&offset=100&id=291047123128
```

- Note that there is query parameters(sort, limit, offset)
- Ref [here](https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/)


**Response**

```http
200 OK
Content-type: application/json
{
    "status": true,
    "data": [
        {
            "id": "0jf0jwe0iajf09i23...",
            "username": "fernando_kim",
            "email": "qwengiewjc@gmeic.com",
            "activity": [
                {
                    "type": "SIGN_IN",
                    "timestamp": 192409124,
                },
                {
                    "type": "ADMIN_NOTE",
                    "content": "I like this guy",
                    "timestamp": 190235912
                },
                ...
            }
            "keychain": [
                {
                    "type": "PASSWORD",
                    "expiration": 0
                },
                {
                    "type": "KAKAO",
                    "content": 2039201,
                    "expiration": 0
                }
            ]
        }
    ]
}
```

- For this week, ignore activity and keychain. We will implement it later
- Just for understanding, 
    - 'activity' is log data of user. If something happens, it will be logged. Note that type of activity is like Admin Note, Sign in, Update, Keychain add etc.
    - 'keychain' is list of credentials the user has. Password data can be stored, Kakao login meta data can be stored, in the future, API keys will be stored.



## Update User

**Request**

```http
PUT /admin/user/{id}
Content-type: application/json
{
    "username": "fernando_kim",
    "email": "qgiwjen.cwof",
    "accessControl": "SYSTEM",
}
```

**Response**

```http
200 OK
Content-type: application/json
{"status": true, "data": {"count": 1}}
```


## Delete User 

**Request**

```http
DELETE /admin/user/{id}
```

**Response**

```http
200 OK
Content-type: application/json
{"status": true, "data": {"count": 1}}
```

# Keychain
- keychain is nested to user
- multiple keychains can be nested user

## Upsert Keychain

**Request**
```http
POST /admin/user/{id}/keychain
Content-type: application/json
{
    "type": "KAKAO",
    "content": "12091482",
    "secret": "",
    "expiration": 0
}
```

- Upsert means this
    - If it doesn't exist, it creates
    - If it exists, it updates
- Unique key is type. If keychain with KAKAO exists, it updates. If PASSWORD doesn't exists, it creates

**Response**
```http
200 OK
Content-type: appication/json
{"status": true, "data": {"count": 1}}
```


## Delete Keychain

**Request**
```http
DELETE /admin/user/{id}/keychain/{keyType}
```

- Allowed {keyType} is "PASSWORD", "KAKAO"

**Response**
```http
{"status": true, "data": {"count": 1}}
```
