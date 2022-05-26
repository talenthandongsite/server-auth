# Data Structure
- This document describes some outline of DB structure
- API structure will be different from the structure described in here. In that case you should map API struct to DB struct.


## User
```json
{
    "id": "string",
    "username": "string",
    "password": "qjweiogj" // this is temporary. Will be removed in future
    "email": "johndoe@acme.com",
    "accessControl": "string enum MASTER|SYSTEM|ADMIN|MEMBER|PENDING|BANNED",
    "activity": [
        {
            "type": "string enum CREATED|UPDATED|ADMIN_NOTE|SIGN_IN|KEYCHAIN_UPSERT|KEYCHAIN_DELETE",
            "content": "string",
            "timestamp": 123094102
        }
    ],
    "keyChain": [
        {
            "type": "string enum PASSWORD|KAKAO",
            "content": "string",
            "secret": "string",
            "expiration": 12949012905
        }
    ]
}
```
