#!/usr/bin/sh

############################################################################
#  PREREQUISITES FOR THIS SCRIPT
# 1. Redistributable-binary(.rpm) should be passed as the first argument
# 2. Login to the cluster should be done prior to running this script
# 3. The cluster should be in a state where it can be used for testing
#
# USAGE:
# ./release-bit-verification.sh redistributable-binary
#
# Example: ./release-bit-verification.sh ~/Downloads/astra-redistributable-2.4.3-1.el8.x86_64.rpm
#
#For erroring out in case of error
set -eo pipefail

shout() {
    echo "--------------------------------$1------------------------------------------"
}
# Check SHASUM for all the binary files and there should be no difference



# Checking for no or invaild arguments
if [ "$#" -lt 1 ]
then
  echo "No arguments supplied"
  exit 1
fi

if [ ! -f ${1} ]; 
then
    echo "Please enter a valid filepath";
    exit 1
fi
#Creating an Temp directory    
WORKING_DIR=$(mktemp -d)
shout "WORKING_DIR=$WORKING_DIR"
export REPO_URL=${REPO_URL:-"https://github\.com/danielpickens/astra.git"}


# Extract from rpm file
rpm2cpio ${1} | cpio -idmvD $WORKING_DIR
pushd $WORKING_DIR/usr/share/astra-redistributable/

# Check sha256sum for all the files
while IFS= read -r line; do
    read -r SHA FILE <<<"$line"
    read -r SHATOCHECK FILE <<<$(sha256sum $FILE)
    if [[ $SHA == $SHATOCHECK ]]; then
        # Print if the file is correct
        printf '%-50s\U0002705\n' $FILE
    fi
done <SHA256_SUM

shout

# Copy binary for testing purpose
OS=$(uname -s)
ARCH=$(uname -m)

if [[ $OS == "Linux" ]]; then
    if [[ $ARCH == "x86_64" ]]; then
        cp ./astra-linux-amd64 astra
        PATH=$(pwd):$PATH
    fi
fi

# Check astra verion and if it is correct
VERSION=$(cat VERSION)
astraVERSIONCHECK=$(astra version)
if [[ "$astraVERSIONCHECK" == *"$VERSION"* ]]; then
    echo "astra binary is installed correctly"
else 
    echo "astra binary is not installed correctly"
    exit 1
fi

#clone repo for testing and checkout release tag
pushd $WORKING_DIR
if [ -d "astra" ]; then
    rm -rf astra
fi
git clone $REPO_URL astra && cd $WORKING_DIR/astra && git checkout "v$VERSION"

#Run tests
make test-e2e

# Cleanup
rm -rf "$WORKING_DIR"


