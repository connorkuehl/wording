## Building

A Go toolchain is all that's needed to build the software.

## Running/operating

### Prerequisites

* [migrate](https://github.com/golang-migrate/migrate). Alternatively, you can manually
  run the SQL statements under migrations/.
* PostgreSQL database.

### Configuring

The application is configured with environment variables or command line flags.
I find it convenient to use `direnv` to load environment variables when I `cd`
into the root of the repo:

```console
$ cat > .envrc <<EOF
export WORDING_DB_DSN='postgres://YOUR_DATABASE_USER:YOUR_PASSWORD@YOUR_DB_HOST/YOUR_DB_NAME'
export WORDING_BIND_ADDR=localhost:8080
export WORDING_BASE_URL=http://localhost
export WORDING_WORD_GEN_SVC='https://random-word-form.herokuapp.com'
EOF
```
