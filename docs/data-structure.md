# Data Structure
- This document describes some outline of DB structure
- API structure will be different from the structure described in here. In that case you should map API struct to DB struct.


## User
```json
{
    "id": "string",                 // id is id of this document, it is automatically generated by mongodb
    "username": "string",           // username is 
    "password": "qjweiogj",         // this is temporary. Will be removed in future
    "email": "johnd@acme.com",      //
    "accessControl": "string",      // enum MASTER|SYSTEM|ADMIN|MEMBER|PENDING|BANNED
    "activity": [                   // activity is 
        {
            "type": "string",           // enum CREATED|UPDATED|ADMIN_NOTE|SIGN_IN|KEYCHAIN_UPSERT|KEYCHAIN_DELETE
            "content": "string",        // content is 
            "timestamp": 123094102      // timestamp is 
        }
    ],
    "keyChain": {                   // keychain is map structure, keys can be anything
        "kakao": {
            "content": "string",        // content is
            "secret": "string",         // 
            "expiration": 12949019      //
        }, 
        "password": {
            "secret": "qjowiergiqwoq"   //
        }
    } 
}
```
