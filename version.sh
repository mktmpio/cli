#!/bin/bash

# Originally copied from https://github.com/fmahnke/shell-semver
# => Copyright (c) 2014 Fritz Mahnke

# Increment a version string using Semantic Versioning (SemVer) terminology.

# Parse command line options.

while getopts ":Mmp" Option
do
  case $Option in
    M ) major=true;;
    m ) minor=true;;
    p ) patch=true;;
  esac
done

shift $(($OPTIND - 1))

# if [ ! -z "$1" ]; then
#   TAG_MSG="-m '$1'"
# fi

version=$(git describe --tags)
version=${version#v}
version=${version%%-*}

# Build array from version string.
a=( ${version//./ } )

# Increment version numbers as requested.

if [ ! -z $major ]
then
  ((a[0]++))
  a[1]=0
  a[2]=0
fi

if [ ! -z $minor ]
then
  ((a[1]++))
  a[2]=0
fi

if [ ! -z $patch ]
then
  ((a[2]++))
fi

next="v${a[0]}.${a[1]}.${a[2]}"
msg=${1:-$next}
git tag -s -a "$next" -m "$msg"
echo $next
