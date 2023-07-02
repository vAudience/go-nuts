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

* v0.1.10 added GetAllStructFieldNames and CopyMatchingFields
* v0.1.9 now filter-tag values for struct field filters use lowercase trimmed string-comparisons and split values by comma beforehand allowing for more complex solutions
* v0.1.8 much better struct filtering
* v0.1.7 added GetStructFieldNamesByTagsValues
* v0.1.6 added FilterStructFields
* v0.1.5 added CopyFields in nuts.structfieldcopy.go. this is meant to be used to filter structs based on access rights.
* v0.1.4 moved to go 1.20 . fixed internal usage of the updated interval. added time coversion helpers for javascript timestamps (Date.now) . added SelectJsonFields and adapted RemoveJsonFields. this will break old usage!
* v0.1.3 updated Interval for the ability to cancel. this changed the call-function type to expecting a returned bool - will break things!
* v0.1.2 added sanitizer func. not really safe! use with caution.
* v0.1.1 updated the ids lib. added the NID() function for more control. updated go dep to 1.19. replaced ioutil.
* v0.1.0 initial isolated version of the nuts package.

## TODO

* nothing so far
