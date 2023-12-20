TestAPI
=======

## TL;DR

```shell
docker run --rm -it -v $PWD/scenarios:/app/scenarios ghcr.io/meysampg/testapi:latest
```

## Scenario Definition ([Sample Scenario File](/scenarios/scenario.yaml-sample))

Do my experiment named `Create Project and Pay` for `2` times that simulate `20` concurrent user on each iteration. I
don't need the output of error on requests.

```yaml
name: Create Project and Pay
iteration: 2
concurrency: 10
verbose: false
```

At first, send a `POST` request to `https://mysite.com/v1/login` with the json
`{'email': 'me@mysite.com', 'password': 'my_very_strong_password'}`  and JSON headers:

```yaml
steps:
  login:
    order: 0
    url: https://mysite.com/v1/login
    verb: post
    pause: 50ms
    headers:
      Accept: 'application/json'
      Content-Type: 'application/json'
    body:
      email: 'me@mysite.com'
      password: 'my_very_strong_password'
```

The response of the previous endpoint is like:

```json
{
  "result": "okay",
  "data": {
    "id": 12,
    "token": "it's_your_bearer_token"
  }
}
```

Wait for 50ms and then use the token on the above response on the next request. On the next request, I want to send
a `POST` request to `https://mysite.com/v1/projects`. It should have `Authorization` and `JSON` headers with a body
like `{"title": "My Project", "slug": "my-project", "related": [1,2,3]}`.

```yaml
  createProject:
    order: 1
    url: https://mysite.com/v1/projects
    verb: post
    pause: 3s
    headers:
      Accept: 'application/json'
      Content-Type: 'application/json'
      Authorization: 'Bearer {{ .login.data.token }}'
    body:
      title: 'My Project'
      slug: 'my-project'
      related: [ 1, 2, 3 ]
```

## Golang Text Template Enabled

Except your first request, you can use [Golang Text Template](https://pkg.go.dev/text/template) everywhere. Suppose the
response of the previous step named `lastUserPayment` is like this:

```json
{
  "id": 11,
  "user": {
    "last_thing": "2021-02-03T04:05:06",
    "payment": {
      "amount": 12345,
      "approved": true
    }
  }
}
```

and you want value of `12345` on your next step. It's possible with `{{ .lastUserPayment.user.payment.amount }}`.