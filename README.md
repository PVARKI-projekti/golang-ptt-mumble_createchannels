# mumble_createchannels

Create channels on mumble server using gumble client library


```bash
go install github.com/PVARKI-projekti/golang-ptt-mumble_createchannels
```


You can also use the `docker` image:

```bash
docker pull pvarkiprojekti/mumble_createchannels
```

Quickstart

```bash
docker run --rm -it -v `pwd`/example.yaml:/channels.yaml:ro pvarkiprojekti/mumble_createchannels:latest --user Superuser --pass "superusers_password" myserver.example.com /channels.yaml
```

Any registered user that is allowed to create channels should work but unless you have exported the certs for said user you can't authenticate as them...

For a full list of options, run with `--help`.

## Development

### pre-commit

Uses pre-commit, you need some basic dependencies (run in this directory)

    pip3 install --user pre-commit detect-secrets
    pre-commit install

As usual using virtualenvs is generally recommended but in this case not strictly mandatory.

Before committing check your work with:

    pre-commit run --all-files ; echo $?

This saves you annoyance of rewriting commit messages when one of the checks fail.

###  Prerequisites

Before you get started, make sure you have installed the following tools::

    $ python3 -m pip install -U cookiecutter>=1.4.0
    $ python3 -m pip install pre-commit bump2version invoke ruamel.yaml halo detect-secrets
    $ go get -u golang.org/x/tools/cmd/goimports
    $ go get -u golang.org/x/lint/golint
    $ go get -u github.com/fzipp/gocyclo/cmd/gocyclo
    $ go get -u github.com/mitchellh/gox  # if you want to test building on different architectures

**Remember**: To be able to excecute the tools downloaded with `go get`,
make sure to include `$GOPATH/bin` in your `$PATH`.
If `echo $GOPATH` does not give you a path make sure to run
(`export GOPATH="$HOME/go"` to set it). In order for your changes to persist,
do not forget to add these to your shells `.bashrc`.


We use bump2version to edit the version numbers in files
```bash
bump2version (major | minor | patch)
git push --follow-tags
```

## Note

This project is still in the alpha stage and should not be considered production ready.
