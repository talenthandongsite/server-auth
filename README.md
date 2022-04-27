> Task of This Week: Implement following five endpoint
>
> - [Health Check](#health-check)
> - [Create User](#create-user)
> - [Read User](#read-user)
> - [Update User](#update-user)
> - [Delete User](#delete-user)
>
> Notes
>
> - When it is done, test it with postman   
> - Use talent-handong Atlas DB (its USERNAME and PASSWORD is in Google Drive)  
> - Make new collection 'user' in the DB, data should be stored in the 'user' collection.  
> - Note that in mongodb all fields will converted to lowercase. (accessControl->accesscontrol) It is expected and just let it do this way.
> - When implementing Create User, field 'passphrase' must be presented. Find a way to validate this.
> - When implementing Update User, if some field is omitted(except for 'passphrase'. When these fields are not presented, API should return error), the field should be set to null or empty value. 
> - Some values must be auto generated. For example, when creating user, 'created' field should be set. Same for 'updated' field. It should be set when the field is updated.
> - When creating or updating user, if field 'accessControl' is not presented, it should be auto populated with 'PENDING'
>
> Recommendation: refer to [mongoDB official package docs](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
>

# Auth Server
Server to provide authentication and authorization of talent-handong.site.

## Background
TBD

### Motivation
TBD

### History
TBD

## Architecture
TBD
 
## Process
TBD


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

## Create User
### Request
```
POST /auth/user
Content-type: application/json
{...userData}
```
> To know what 'userData' is, refer to [User](#user)
### Response
```
200 OK
Content-type: application/json
{"status": true, "data": {"count": 1}}
```
> The field 'count' means the number of inserted data

## Read User
### Request
```
GET /auth/user
```
### Response
```
200 OK
Content-type: application/json
{"status": true, "data": [...userData]}
```
> To know what 'userData' is, refer to [User](#user)

## Update User
### Request
```
PUT /auth/user/{id}
Content-type: application/json
{...userData}
```
> To know what 'userData' is, refer to [User](#user)

### Response
```
200 OK
Content-type: application/json
{"status": true, "data": {"count": 1}}
```
> The field 'count' means the number of updated data

## Delete User 
### Request
```
DELETE /auth/user/{id}
```
### Response
```
200 OK
Content-type: application/json
{"status": true, "data": {"count": 1}}
```
> The field 'count' means the number of deleted data

# Reference

## User 
<table style="width: 100%">
    <thead>
        <th>Name</th>
        <th>Type</th>
        <th>Description</th>
        <th>Note</th>
    </thead>
    <tbody>
        <tr>
            <td>id</td>
            <td>string</td>
            <td>The key of this user data</td>
            <td>Note that this field should be omitted in 'Create' and 'Update' operation</td>
        </tr>
        <tr>
            <td>externalId</td>
            <td>string</td>
            <td>Key for external system. From this point, it is kakaotalk id</td>
            <td></td>
        </tr>
        <tr>
            <td>username</td>
            <td>string</td>
            <td>Nickname of user(for my case: '김페르난도')</td>
            <td></td>
        </tr>
        <tr>
            <td>passphrase</td>
            <td>string</td>
            <td>Passphrase for authorization</td>
            <td>Required</td>
        </tr>
        <tr>
            <td>accessControl</td>
            <td>string</td>
            <td>Access control level of this user</td>
            <td>Default is 'PENDING'. Refer to <a href="#access-control">Access Control</a></td>
        </tr>
        <tr>
            <td>created</td>
            <td>number(int64)</td>
            <td>Created timestamp of this user</td>
            <td>This value should be unix time value. In Go, use (Time).UnixMilli()</td>
        </tr>
        <tr>
            <td>updated</td>
            <td>number(int64)</td>
            <td>Updated timestamp of this user</td>
            <td>This value should be unix time value. In Go, use (Time).UnixMilli()</td>
        </tr>
        <tr>
            <td>lastAccess</td>
            <td>number(int64)</td>
            <td>Timestamp for last access</td>
            <td>This value should be unix time value. In Go, use (Time).UnixMilli()</td>
        </tr>
        <tr>
            <td>adminNote</td>
            <td>string</td>
            <td>Note for admin.</td>
            <td></td>
        </tr>
    </tbody>
</table>

## Access Control
<table style="width: 100%">
    <thead>
        <th>Key</th>
        <th>Description</th>
    </thead>
    <tbody>
        <tr>
            <td>MASTER</td>
            <td>Unlimited access to all features</td>
        </tr>
        <tr>
            <td>SYSTEM</td>
            <td>system user</td>
        </tr>
        <tr>
            <td>ADMIN</td>
            <td>Personnel who should access to data for maintenance</td>
        </tr>
        <tr>
            <td>MEMBER</td>
            <td>Ordinary member. Only access to frontend view</td>
        </tr>
        <tr>
            <td>PENDING</td>
            <td>Default value. This user access to nothing but should be updated to MEMBER</td>
        </tr>
        <tr>
            <td>BANNED</td>
            <td>Banned by admin. This user access to nothing</td>
        </tr>
    </tbody>
</table>
