## This is used in openshiftci-presubmit-all-tests.sh, used by Interop tests.
## Since these interop tests are managed by a different team, we intentionally want to use a stable Devfile registry.
## But our own internal tests make use of the staging Devfile registry.
# NOTE: Since #7154, tests no longer depend on an external Devfile Registry. They are now able to start an HTTP Server emulating the Devfile Registry API.
#export DEVFILE_REGISTRY=https://devfile-registry-ci-devfile-registry.astra-test-kubernete-449701-49529fc6e6a4a9fe7ebba9a3db5b55c4-0000.us-east.containers.appdomain.cloud/
