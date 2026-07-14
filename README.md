RSS Feed aggre`gator`

Built with:
- Golang
    https://go.dev/doc/install
- Postgres
    https://www.postgresql.org/download/

Need to install everything above to be able to use this program.

Also built with these but they don't need to be installed to use gator:
- SQLC
    https://docs.sqlc.dev/en/latest/overview/install.html
- GOOSE
    https://github.com/pressly/goose

1. Set up

Run:

```bash
go install github.com/yourusername/gator@latest
```

Set up these required files before running the program

- ~/.gatorconfig.json

!!! Don't forget:
Replace user:pass with your own PostgreSQL username and pass word before running the below command

```bash
cat > ~/.gatorconfig.json <<'EOF'
{"db_url":"postgres://user:pass@localhost:5432/dbname?sslmode=disable"}
EOF
```
-

Do this to use command gator globally:

cd to your aggreGATOR folder

```bash
cd /aggreGATOR
```

Then, run these commands:
```bash
mkdir -p "$(go env GOPATH)/bin"
```
```bash
go build -o "$(go env GOPATH)/bin/gator" .
```

Then, check to make sure $(go env GOPATH)/bin is in your PATH:
```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

2. Basic commands

Any command usage:

- in the terminal type out: 

```bash
gator
```

```bash
gator <command> <arguments>
```

Command list:
- register: register username

```bash
gator register <user>
```

    Ex: register gatoruser

- login: login with username

```bash
gator login <user>
```

    Ex: login gatoruser

- agg: show feed every <time>

```bash
gator agg <time>
```

    Ex: agg 2h
        ---> This shows feed every 2 hours
    valid time: s, m, h

- addFeed: requires the user to be logged in, add a feed to database 

```bash
gator addfeed "<name>" "<url>"
```

    Ex: addfeed "newfeed" "feedURL"

- feeds: show current feeds in database

```bash
gator feeds
```

- follow: requires the user to be logged in, follow a feed

```bash
gator follow "<url>"
```

- following: requires the user to be logged in, shows all feeds the current user is following

```bash
gator following
```
- unfollow: requires the user to be logged in, unfollow a feed

```bash
gator unfollow "<url>"
```
- browse: show default 2 posts, can add an argument to change how many posts shown

```bash
gator browse
```

        ---> Shows 2 posts

```bash
gator browse 10
```

        ---> Shows 10 posts