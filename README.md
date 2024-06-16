# PRJ cli tool to quickly jump to your projects

## Install

### Prerequisites
have `go` installed and available in `$PATH` on your machine

Create `~/.prjrc` file in your home directory containing a single variable that contains a Glob pattern for your projects.
```bash
PROJECTS_HOME=$HOME/repos/*/*
```

build the tool with make
```bash
make install
```
