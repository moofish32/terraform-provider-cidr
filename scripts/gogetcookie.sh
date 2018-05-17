#!/bin/bash
set +o history
touch ~/.gitcookies
chmod 0600 ~/.gitcookies

git config --global http.cookiefile ~/.gitcookies

tr , \\t <<\__END__ >>~/.gitcookies
go.googlesource.com,FALSE,/,TRUE,2147483647,o,git-moofish32.gmail.com=1/ziSvpbUBIqftGXePPhyKmpqoFrlEyW7kVLFCX_Op-GEovUPxec4eJOlQ_OVQnYi1
go-review.googlesource.com,FALSE,/,TRUE,2147483647,o,git-moofish32.gmail.com=1/ziSvpbUBIqftGXePPhyKmpqoFrlEyW7kVLFCX_Op-GEovUPxec4eJOlQ_OVQnYi1
__END__
set -o history
