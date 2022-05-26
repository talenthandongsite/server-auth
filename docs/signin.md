
# Sign In

**Request**
```http
POST /signin
Content-type: application/json
{
    "username": "fernando_kim",
    "password": "secretText12!"
}
```

- for the last example, we used keyword 'passphrase'. Change it to 'password'
- note that it's endpoint is '/signin'

**Response**
```http
200 OK
Content-type: application/json
{
    "status": true,
    "data": {
        "token": "Bearer qpigjwepijgqjwegpijqpwe...",
        "exp": 1923012492
    }
}
```

- When request comes in, username and password should be checked in database
- 'token is jwt token. search for it
- In token claim these information should be included 
    - id (database id of this user)
    - username
    - accessLevel
- The keyword 'Bearer' and single space character should be inserted before token string.