# mktmpio CLI

## Installation

1. Sign up at https://mktmp.io/
2. Download the binaries for your platform from the [latest release](https://github.com/mktmpio/cli/releases/latest).
3. Extract the binary.
4. Run the binary from the location you extracted it to or add it to your path.

## Usage

Once installed and your `~/.mktmpio.yml` config contains [your mktmpio auth
token](https://mktmp.io/me) you can create an instance of any of the supported services:

    $ mktmpio $TYPE

### Examples

#### Redis

```
$ mktmpio redis
i.mktmp.io:32806> SCAN 0
1) "0"
2) (empty list or set)
i.mktmp.io:32806>exit
Instance e19b07bca586 terminated.
$
```

#### MySQL

```
$ mktmpio mysql
mysql> select 1, 2;
+---+---+
| 1 | 2 |
+---+---+
| 1 | 2 |
+---+---+
1 row in set (0.07 sec)

mysql> exit
Bye
Instance 3b9f136893da terminated.
$
```

## Development

### Build

1. `git clone git@github.com:mktmpio/cli.git $GOPATH/github.com/mktmpio/cli`
2. `cd $GOPATH/github.com/mktmpio/cli`
3. `go get -u -v -t ./...`
4. `go build -v`
5. You should now be able to run `mktmpio` as `./cli`

### Release Procedure

1. create an annotated release tag, following SemVer
  * major: `./version.sh -M`
  * minor: `./version.sh -m`
  * patch: `./version.sh -p`
2. push tag to github
  * tags trigger release builds on Travis
  * can be manually run with `make release`
  * release binaries are uploaded to GitHub release automatically

## Legal

This software is &copy; 2015 Datajin Technologies, Inc. and Open Source under
the terms of Artistic License 2.0 that can be found in the LICENSE file.

Use of the mktmpio service is subject to the
[mktmpio Privacy Policy](https://mktmp.io/privacy-policy) and
[mktmpio Terms of Service policy](https://mktmp.io/terms-of-service).

---
&copy; 2015 Datajin Technologies, Inc.
