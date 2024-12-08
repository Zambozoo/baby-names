# Baby Names
A quick little command line project to help couples pick a name for their baby.

Couples are prompted to set up accounts and asked to like or dislike names as they appear.
Any matches are shown on login and when they are initially matched.

## Running
Use [taskfiles](https://taskfile.dev/) for an easy start up:
```bash
# sets USER_DB=db.bolt NAMES_FILE=names/boys.txt
task go:boy
task go:girl
```

The environment varaibles are optional if running without taskfiles:
```bash
go run ./cmd/run/...
```

## Clearing Users
To clear a couple from the database:
```bash
# sets USER_DB=db.bolt
task go:clear
```

The environment variable is optional if running without taskfiles:
```bash
go run ./cmd/reset/...
```