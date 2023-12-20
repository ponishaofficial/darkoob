Darkoob
=======

> [!NOTE]  
> ![Woodpecker-habitat](https://github.com/ponishaofficial/darkoob/assets/1416085/a3bd1283-90ca-4a14-92fd-290725077eb3)
> Image from [WorldAnimalFoundation](https://worldanimalfoundation.org/advocate/wild-animals/params/post/1292094/woodpeckers).
> Darkoob (in Persian, [pronunciation](https://forvo.com/word/%D8%AF%D8%A7%D8%B1%DA%A9%D9%88%D8%A8/)) or Woodpeckers are part of the bird family Picidae, which also includes the piculets, wrynecks and sapsuckers. Members of this family are found worldwide, except for Australia, New Guinea, New Zealand, Madagascar, and the extreme polar regions.
> __ [Wikipedia](https://en.wikipedia.org/wiki/Woodpecker).
## TL;DR

```shell
git clone https://github.com/ponishaofficial/darkoob
cd darkoob
cat scenarios/scenario.yaml-sample # view the sample
cp scenarios/scenario.yaml-sample scenarios/first-scenario.yaml # get a copy from sample
vim scenarios/first-scenario.yaml # edit based on your need
docker run --rm -it -v $PWD/scenarios:/app/scenarios ghcr.io/ponishaofficial/darkoob:latest # run your tests
```

## Scenario File Crash Course

```yaml
name: String -Name of the scenario to show on final report
iteration: Integer -How many times this scenario should run? Use 0 to disable scenario
concurrency: Integer -How many concurrent users should be simulated?
silent: Boolean -Show progress of steps or not
verbose: Boolean -Show error details or not
steps: # Definition of steps, for example:
  login: # Name of the step. It must be in camelCase PascalCase format (S*#t! Again I'm doubtful about their names 🤔)
    order: Integer -Shows the order of the step on running. Remove it to send requests in a nondeterministic order.
    url: String -URL
    verb: String -get, post, delete, patch or put
    pause: String -How much should pause after the request of this step? Must be in X😜 format which 😜 is ms (millisecond), s (second), m (minute), and h (hour), a.e. 12m13s is 12 minutes and 13 seconds.
    headers: # headers of the request, keys, and values must be a string
      HeaderName: HeaderValue # for example:
      Content-Type: 'application/json'
    body: # body of the request that will be converted into JSON
      key: value # for example:
      email: 'my@awsome.email'
      list_of_ids: [ 1, 2, 3 ]
      or_another_object:
        hello: 'my_friend'
```

## Scenario Definition ([Sample Scenario File](/scenarios/scenario.yaml-sample))

Do my experiment named `Create Project and Pay` for `2` times that simulate `20` concurrent users on each iteration. I
don't need the output of error on requests.

```yaml
name: Create Project and Pay
iteration: 2
concurrency: 20
verbose: false
```

First, send a `POST` request to `https://mysite.com/v1/login` with the JSON
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

Except for your first request, you can use [Golang Text Template](https://pkg.go.dev/text/template) everywhere. Suppose the
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

and you want a value of `12345` on your next step. It's possible with `{{ .lastUserPayment.user.payment.amount }}`.

### Supported Inline Functions
- `{{ split .fieldName "separator" indexToReturn }}`, for example if email on data object is `my@awesome.email`, then
  `{{ split .data.email "@" 1 }}` returns `my`.

You're welcome to add your desired functions to [funcs.go](/utils/funcs.go) and create a PR.

***
Darkoob is made with ❤️ for you in [Ponisha](https://ponisha.ir).