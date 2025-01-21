#!/usr/bin/env bash

set -e

echo "Preping rpm"
scripts/rpm-prepare.sh

echo "Building rpm"
scripts/rpm-local-build.sh

rm -rf dist/rpmtest
mkdir -p dist/rpmtest/{astra,redistributable}

echo "Validating astra rpm"
rpm2cpio dist/rpmbuild/RPMS/x86_64/`ls dist/rpmbuild/RPMS/x86_64/ | grep -v redistributable` > dist/rpmtest/astra/astra.cpio
pushd dist/rpmtest/astra
cpio -idv < astra.cpio
ls ./usr/bin | grep astra
./usr/bin/astra version
popd

RL="astra-darwin-amd64 astra-darwin-arm64 astra-linux-ppc64le astra-linux-arm64 astra-windows-amd64.exe astra-linux-amd64 astra-linux-s390x"
echo "Validating astra-redistributable rpm"
rpm2cpio dist/rpmbuild/RPMS/x86_64/`ls dist/rpmbuild/RPMS/x86_64/ | grep redistributable` > dist/rpmtest/redistributable/astra-redistribuable.cpio
pushd dist/rpmtest/redistributable
cpio -idv < astra-redistribuable.cpio
for i in $RL; do
	ls ./usr/share/astra-redistributable | grep $i
done
./usr/share/astra-redistributable/astra-linux-amd64 version
popd
