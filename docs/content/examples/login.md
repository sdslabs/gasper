# Login

This example shows how to login into the Gasper ecosystem and obtain a JSON Web Token

!!!info
    The JSON Web Token obtained will be used in other requests

!!!warning "Prerequisites"
    You have [Master](/configurations/master/) up and running

```bash
$ curl -X POST \
  http://localhost:3000/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "anish.mukherjee1996@gmail.com",
    "password": "alphadose"
  }'

{
    "code": 200,
    "expire": "2019-12-04T22:05:41+05:30",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJhbHBoYWRvc2VAZ21haWwuY29tIiwiZXhwIjoxNTc1NDc3MzQxLCJvcmlnX2lhdCI6MTU3NTQ3Mzc0MSwidXNlcm5hbWUiOiJhbHBoYWRvc2UifQ.Io0txryVH8zR6JfZ0iey86474oZl8gNwo4HjKgZl2s8"
}
```

The **token** field in the above response holds the required JSON Web Token
