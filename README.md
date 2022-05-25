> Task of This Week
> - Get Used to [Github Project](https://github.com/talenthandongsite/server-auth/projects/1)
> - Implement sign in
> - Implement keychain
>

# Auth Server
Server to provide authentication and authorization of talent-handong.site.

## Motivation
Some HGT tools includes confidential information, so it needs some kind of access control over resources.

Also HGT is Kakao Talk based community, and Kakao Talk provides [OAuth2.0](https://datatracker.ietf.org/doc/html/rfc6749) compatible [login servic](https://developers.kakao.com/docs/latest/ko/kakaologin/rest-api).

## Restrictions
- HGT is closed community and new registration must be authorized by administrator
- [Kakao Talk sign in](https://developers.kakao.com/docs/latest/ko/kakaologin/rest-api) should be implemented


## Business rules
1. User can sign in with Account/PW as well as 3rd Party Authorization such as OAUTH2 or OpenIDConnect  
    - this is because there are multiple use cases
    - the first one is member, they use kakao login
    - the second one is system user, which is not an actual person, so that it doesn't have any kakao account

2. Token would be provided after login, these should expire 30days after. System user's token will never expires.

3. With token, user can access resource server's API.

4. token has Access Level. AL will be used to restrict confidential data from user. For example, administration API of this server will only accessable to user who calls API with ADMIN token.


## Architecture and Process
![arch](/assets/img/architecture.png)

### Sign in Process
1. From client, start Kakao login process. When process ends, Kakao talk issues a token
2. Client send token to server
3. Server validate token and get information in it
4. Update information to database

### Administration Process
1. Client calls API
2. Do crud on database
3. Return result to client

## What type of Signin 

# REST API

## Health Check
Simple health check API.

### Request
```
GET /
```
### Response
```
200 OK
Content-type: application/json
[HGT Auth Server] Server is good.
```

## Administration
[User Administration API](/docs/admin-user.md)

## Sign In
[Sign In API](/docs/signin.md)
