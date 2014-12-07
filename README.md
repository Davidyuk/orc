ORC
===

Course work.

**ORC** helps to simplify the interaction of the participants with
the organizers of the various competitions, events, tournaments, seminars,
conferences, etc.

Installation
------------------------

### Requirements

[Go][1] v1.2 or higher.

[PostgreSQL][2] v9.3 or higher.

### Getting Source Files

Install [git][3]. Get a copy of repository:

    $ git clone git@github.com:GuraYaroslava/orc.git

### Install Packages

- [package pq][4] ([godoc](http://godoc.org/github.com/lib/pq))

        $ go get github.com/lib/pq

- [package securecookie][5] ([godoc](http://godoc.org/github.com/gorilla/securecookie))

        $ go get github.com/gorilla/securecookie

### Import Dump File

Use PSQL Console:

    orc=# \i path/to/dumpfile

### Running

    $ go build && orc.exe

[1]: https://golang.org
[2]: http://www.postgresql.org
[3]: http://git-scm.com
[4]: https://github.com/lib/pq
[5]: http://www.gorillatoolkit.org/pkg/securecookie
