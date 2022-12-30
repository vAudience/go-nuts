# go-nuts

n-utils for go private repository

## HOW-TO ACCESS

reference is here: <https://go.dev/doc/faq#git_https>

1. create a Github Personal Access Token <https://github.com/settings/tokens/new> with:
  1.1. repo full (read+write)
  1.2. write+read packages
  1.3. projects read access
2. update ```$HOME/.netrc``` and add

````bash
machine github.com login USERNAME password PERSONAL-ACCESS-TOKEN
````

## VERSIONS

* 0.1.2 added sanitizer func. not really safe! use with caution.
* 0.1.1 updated the ids lib. added the NID() function for more control. updated go dep to 1.19. replaced ioutil.
* 0.1.0 initial isolated version of the nuts package.

## TODO

* nothing so far
