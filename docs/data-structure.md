# Data Structure
- This document describes some outline of DB structure
- API structure will be different from the structure described in here. In that case you should map API struct to DB struct.


## User
```json
{
    "id": "string",
    "username": "string",
    "passphrase": "string",
    "accessControl": "MASTER|SYSTEM|ADMIN|MEMBER|PENDING|BANNED",
    "created": 1293192,
    "updated": 1924091,
    "activity": [
        {
            "type": "ADMIN_NOTE|SIGNIN",
            "content": "string",
            "timestamp": 123094102
        }
    ],
    "keyChain": [
        {
            "type": "PASSPHRASE|KAKAO|API_KEY",
            "content": "string",
            "secret": "string",
            "expiration": 12949012905
        }
    ]
}
```
