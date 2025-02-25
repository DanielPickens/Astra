#!/bin/bash

# this script assumes that runs on linux

BIN_DIR="./dist/bin"
RELEASE_DIR="./dist/release"

mkdir -p $RELEASE_DIR

# if this is run on travis make sure that binary was build with corrent version
if [[ -n $TRAVIS_TAG ]]; then
    echo "Checking if astra version was set to the same version as current tag"
    # use sed to get only semver part
    bin_version=$(${BIN_DIR}/linux-amd64/astra version --client | head -1 | sed "s/^astra \(.*\) (.*)$/\1/")
    if [ "$TRAVIS_TAG" == "${bin_version}" ]; then
        echo "OK: astra version output is matching current tag"
    else
        echo "ERR: TRAVIS_TAG ($TRAVIS_TAG) is not matching 'astra version' (v${bin_version})"
        exit 1
    fi
fi

# gziped binaries
for arch in `ls -1 $BIN_DIR/`;do
    suffix=""
    if [[ $arch == windows-* ]]; then
        suffix=".exe"
    fi
    source_file=$BIN_DIR/$arch/astra$suffix
    source_dir=$BIN_DIR/$arch
    source_filename=astra$suffix
    target_file=$RELEASE_DIR/astra-$arch$suffix

    # Create a tar.gz of the binary
    if [[ $suffix == .exe ]]; then
        echo "zipping binary $source_file as $target_file.zip"
        zip -9 -y -r -q $target_file.zip $source_dir/$source_filename
    else
        echo "gzipping binary $source_file as $target_file.tar.gz"
        tar -czvf $target_file.tar.gz --directory=$source_dir $source_filename
    fi

    # Move binaries to the release directory as well
    echo "copying binary $source_file to release directory"
    cp $source_file $target_file
done

function release_sha() {
    echo "generating SHA256_SUM for release packages"
    release_dir_files=`find $RELEASE_DIR -maxdepth 1 ! -name SHA256_SUM -type f -printf "%f\n"`
    for filename in $release_dir_files; do
        sha_sum=`sha256sum $RELEASE_DIR/${filename}|awk '{ print $1 }'`; echo $sha_sum  $filename;
    done > ${RELEASE_DIR}/SHA256_SUM
    echo "The SHA256 SUM for the release packages are:"
    cat ${RELEASE_DIR}/SHA256_SUM
}

release_sha