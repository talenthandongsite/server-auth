> Task of This Week
> - Get Used to [Github Project](https://github.com/talenthandongsite/server-auth/projects/1)
> - Implement sign in (both)
> - Change Data structure (both)
> - Implement keychain (투달)
> - Implement query parameter for Read User (y3x)
>

# Auth Server
Server to provide authentication and authorization of talent-handong.site.

## Motivation
As an investment group, HGT holds some confidential information from various sources. These data are usually shared through messengers and not shared outside. However, after some time goes, HGT started to automated its manual works and transforming data to web service like things. First of all, HGT had to block access from public, and also, need to restrict access to some of data only to administrative members. Also, it is expected that some system users are required in the future.

Since we utilize certain messenger([KakaoTalk](https://www.kakaocorp.com/page/service/service/KakaoTalk)), social login with this messenger is required. Also with system users, default login should be required. It means that authentication service should provide multi-method login.

HGT service aims for 2 things, 1) many features, 2) easy to implement each feature. With these aims, we choose to use micro-service kind of paradigm on our service. That is the reason why separated authentication server is required. As [Netflix](https://netflixtechblog.com/edge-authentication-and-token-agnostic-identity-propagation-514e47e0b602) described, managing identities from different services are complex and inefficient.

OAUTH2 and JWT based identity is good choice for today. OAUTH2 describes good method for issuing identity from auth server, and JWT describes for easy and secure identity token.

## Goals 
- Access Levels for user
- Multi-Method login including [Kakao Talk sign in](https://developers.kakao.com/docs/latest/ko/kakaologin/rest-api) should be implemented
    - For the initial dev, we should implement 1) Kakao, 2) ID/PW
- Login should provides JWT token. This token should be used in all services.


## Architecture and Process
![arch](/assets/img/architecture.png)

### Components
- Auth Server: the server we should implement
- Frontend: we will have frontend that has login form(ID/PW) or social login buttons
- DB: BSON based database
- External Login Page: OAUTH2 or OpenIDConnect based web login page. It will redirect to frontend after login.
- External Login Server: server that provide identity verification(check expiration or invalidity when send token to it

### KAKAO Sign In Process(RED)
1. From client, start Kakao login process. When process ends, Kakao talk issues a token
2. Client send token to server
3. Server validate token and get information in it
4. Update information to database
5. Forge JWT Token
6. Issue Token

### ID/PW Sign In Process(RED, but omit some step)
2. Send ID/PW data
4. Verify ID/PW data
5. Forge JWT Token
6. Issue Token

### Administration Process(ORANGE)
1. Client calls API
2. Do crud on database
3. Return result to client

# Data Structure
[Data Structure](/docs/data-structure.md)

# REST API
[REST API](/docs/endpoint.md)

