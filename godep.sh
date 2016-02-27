#!/bin/sh -x
rm -f Godeps/Godeps.json
godep save -r -t .
