#!/usr/bin/bash

ROLLBACK_VERSION="$1"
RANCHER="ranger/ranger"
ROLLBACK_IMAGE_TAG="${RANCHER}:$ROLLBACK_VERSION"
CONTAINER_ID=`docker ps | awk 'NR > 1 {print $1}'`
DATA_CONTAINER="ranger-data"
VARLIB="/var/lib/ranger"
BACKUP="/backup/${DATA_CONTAINER}-${ROLLBACK_VERSION}.tar.gz"

rollbackRanger() {
    echo -e "\nStopping Ranger..."
    docker stop "${CONTAINER_ID}"

    echo -e "\nPulling old Ranger image..."
    docker pull "${ROLLBACK_IMAGE_TAG}"

    echo -e "\nReplacing data in ${DATA_CONTAINER} with the data in ${BACKUP}..."
    docker run --volumes-from "${DATA_CONTAINER}" -v ${PWD}:/backup busybox sh -c "rm ${VARLIB}/* -rf && tar zxvf ${BACKUP}"

    echo -e "\nStarting Ranger..."
    docker run -d --volumes-from "${DATA_CONTAINER}" --restart=unless-stopped \
                                                     -p 80:80 -p 443:443 \
                                                     --privileged "${ROLLBACK_IMAGE_TAG}"
}

usage() {
	cat << EOF

$(basename "$0")

This script will rollback Ranger API Server using Docker. This script assumes the following:

    * Ranger is running in a Docker container
    * Docker is installed and script user is in the docker group
    * The upgrade.sh script has been run before this one

When running the script, specify the version of Ranger to rollback to, prefixed with a leading 'v'.

USAGE: % ./$(basename "$0") [options]

OPTIONS:
	-h	-> Usage

EXAMPLES OF USAGE:

* Run script
	
	$ ./$(basename "$0") v<version to rollback to>

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
    rollbackRanger
}

Main "$@"
