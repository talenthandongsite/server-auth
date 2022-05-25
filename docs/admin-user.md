User is

## Create User
### Request
```
POST /admin/user
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
GET /admin/user
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
PUT /admin/user/{id}
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
DELETE /admin/user/{id}
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
            <td>Default value. This user access to nothing but would be updated to MEMBER</td>
        </tr>
        <tr>
            <td>BANNED</td>
            <td>Banned by admin. This user access to nothing</td>
        </tr>
    </tbody>
</table>
