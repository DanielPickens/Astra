#!/bin/bash

set -e

asciioutput() {
  OUTPUT=`${@}`
  ALLOUTPUT="== ${@}

[source,json]
----
${OUTPUT}
----
"
  echo ${ALLOUTPUT} > test.adoc
}

tmpdir=`mktemp -d`
cd $tmpdir
git clone https://github.com/openshift/nodejs-ex
cd nodejs-ex

# Commands that don't have json support
# app delete
# catalog describe
# catalog search
# component create
# component delete
# component link
# component log
# component push
# component unlink
# component update
# component watch
# config set
# config unset
# config view
# debug port-forward 
# preference set
# preference unset
# preference view
# service create
# service delete
# storage delete
# url create
# url delete
# login
# logout
# utils *
# version


# Alphabetical order for json output...

# Preliminary?
astra project delete foobar -f || true
sleep 5
astra project create foobar
sleep 5
astra create nodejs
astra push

# app
asciioutput astra app describe app -o json
astra app list -o json

# catalog
astra catalog list components -o json
astra catalog list services -o json

# component
astra component delete -o json
astra component push 

# project
astra project create foobar -o json
astra project delete foobar -o json
astra project list -o json

# service

## preliminary
astra service create mongodb-persistent mongodb --plan default --wait -p DATABASE_SERVICE_NAME=mongodb -p MEMORY_LIMIT=512Mi -p MONGODB_DATABASE=sampledb -p VOLUME_CAPACITY=1Gi
astra service list -o json

# storage
astra storage create mystorage --path=/opt/app-root/src/storage/ --size=1Gi -o json
astra storage list -o json
astra storage delete

# url
astra url create myurl
astra url list -o json
astra url delete myurl
