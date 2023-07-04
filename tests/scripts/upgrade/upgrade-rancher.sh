#!/usr/bin/bash

OLD_VERSION="$1"
NEW_VERSION="$2"
RANCHER="ranger/ranger"
OLD_IMAGE_TAG="${RANCHER}:${OLD_VERSION}"
NEW_IMAGE_TAG="${RANCHER}:${NEW_VERSION}"
CONTAINER_ID=`docker ps | awk 'NR > 1 {print $1}'`
DATA_CONTAINER="ranger-data"
VARLIB="/var/lib/ranger"
BACKUP="/backup/${DATA_CONTAINER}-${OLD_VERSION}.tar.gz"

upgradeRanger() {
    echo -e "\nStopping Ranger and creating a data container..."
    docker stop "${CONTAINER_ID}"
    docker create --volumes-from "${CONTAINER_ID}" --name "${DATA_CONTAINER}" "${OLD_IMAGE_TAG}"

    echo -e "\nCreating a backup tarball..."
    docker run --volumes-from "${DATA_CONTAINER}" -v "${PWD}:/backup" --rm busybox tar zcvf "${BACKUP}" "${VARLIB}"

    echo -e "\nPulling new Ranger image..."
    docker pull "${NEW_IMAGE_TAG}"

    echo -e "\nStarting Ranger..."
    docker run -d --volumes-from "${DATA_CONTAINER}" --restart=unless-stopped \
                                                     -p 80:80 -p 443:443 \
                                                     --privileged "${NEW_IMAGE_TAG}"
}

usage() {
	cat << EOF

$(basename "$0")

This script will upgrade Ranger API Server using Docker. This script assumes the following:

    * Ranger is running in a Docker container
    * Docker is installed and script user is in the docker group

When running the script, specify the current Ranger version and the upgraded Ranger version.

Both need to be prefixed with a leading 'v'.

USAGE: % ./$(basename "$0") [options]

OPTIONS:
	-h	-> Usage

EXAMPLES OF USAGE:

* Run script
	
	$ ./$(basename "$0") v<current Ranger version> v<upgraded Ranger version>

EOF
}

while getopts "h" opt; do
	case ${opt} in
		h)
			usage
			exit 0;;
    esac
done

Main() {
    upgradeRanger
}

Main "$@"
