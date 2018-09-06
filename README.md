# chainstack-coding-assignment
Chainstack quota management API server.

## Technologies choices

The API server is built in Golang. Postgres for datastore.
Gin is used as the web engine.



These are just personal preferences, my experiences saw these pieces work well for most cases while also give fast development pace.



## Authentication
1. To authenticate a user, a form of email-password is sent in plain text to the server. 
2. The server then generate a random `salt` and compute the password hash as followed:

```
hashedPasswordToStoreInDB := bcrypt(randomizeSalt + userInputPassword)
```
- I use `bcrypt` instead of `md5`, `sha1` or even `sha512`. While `md5` and `sha1` is proved to have collisions, `sha512` is a fast function that might give a advantage to attackers who has powerful computers.
- Randomized `salt` is also stored in the db along with `email` and `hashPassword`.
 
3. Comparing `hashedPassword` and `passwordInDB` would tell us whether authentication is successfully or not.

4. Upon successfully authenticate, the server generates a access token and tells the browser to set cookie to this value by responding with a `Set-Cookie` header.

5. The browser would then use that cookie for later requests. For manual calls to the API, we would have to use the header `X-Access-Token`.



## Authorization

I have middlewares to check:
- whether user is login or not
- who is him
- is he an admin



All these are done by decoding the access token to get user email. The middlewares then make a query to the database to get the full user object with that email.



## Quota management


I have a separate table to manage user quota, called `user_quotas`.


```
user_quotas(
	id INTEGER,
	user_id INTEGER,
	quota INTEGER,
	current_quota_left INTEGER
)
```


All write operations that interact with this table is wrapped inside a Postgresql transaction. Specifically:
- Newly created users dont have a record in this table yet. (infinite quota)
- When quota is updated, `current_quota_left` is reset to `quota`
- If a user want to create a resource, this condition `current_quota_left > 0` must satisfies
- Once created a resource, `current_quota_left` is decreased by 1
- Once deleted a resource, `current_quota_left` is increase by 1, until it equals `quota`

## Development

I wrote a `Makefile` to fire up the server locally for development. 
- `make dev` fire up the server locally
- `make build` build the binary
- `make start` start the server in background using `nohup`
- `make restart` gracefully restart the server with the newly built binary

I also set up different configs for local and production environment.
- `config.yaml` is ignored by `.gitignore`, it is the local config file.
- `config.yaml.example` is the sample config file. When deploying, I would look into this file and see what does it need, all I have to do is fill in the blank.


To manage database migrations, I use `goose`.
- `goose up` upgrades migration version to the latest one.
- `goose down` jumps down 1 migration version
- `goose create sample_migration sql` creates a sql migration file with timestamp.
- `goose status` shows the current migration version of my database.

## Tests
To be completely honest, I don't usually write tests in my daily workflow. That's something I wish to change.

Although, in this project, I did try to write 1 test, the login handler test. After that, I feel like I'm shooting at my own feet because I have something designed wrongly so it isn'y easy to write test for it.

## Deployment
The project is deployed to my personal VPS. The API can be found at `chainstack.laptrinhviendeptrai.xyz`.
