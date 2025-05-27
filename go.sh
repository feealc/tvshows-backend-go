#! /bin/bash

PROGRAM_NAME=$(basename $0)
FLAG_BUILD_RUN=0
FLAG_TESTS=0
FLAG_SET_DEBUG=0

# =================================================================================================

usage()
{
    echo
    echo -e "usage:"
    echo -e "\t ${PROGRAM_NAME} flags2"
    echo
    echo -e "\t flags"
    echo -e "\t [-d] set debug flag true"
    echo -e "\t [-r] go run"
    echo -e "\t [-t] go test"
    echo
    exit 1
}

exec_common()
{
    unset CLOUD_NAME
}

exec_build_run()
{
    exec_common

    export DB_HOST="localhost"
    export DB_PORT="5432"
    export DB_USER="root"
    export DB_PASS="root"
    export DB_NAME="root"

    # Go
    go run main.go
}

exec_tests()
{
    exec_common

    export DB_HOST="localhost"
    export DB_PORT="5433"
    export DB_USER="root"
    export DB_PASS="root"
    export DB_NAME="root"

    # Go
    clear
    go test -v ./tests -count=1
}

# =================================================================================================

if [[ ${1} == "" ]]; then
    usage
fi

while getopts "drt" flag; do
    case "${flag}" in
        d)
            FLAG_SET_DEBUG=1
            ;;

        r)
            FLAG_BUILD_RUN=1
            ;;

        t)
            FLAG_TESTS=1
            ;;

        *)
            usage
            ;;
    esac
done

if [[ ${FLAG_SET_DEBUG} -eq 1 ]]; then
    export DEBUG="true"
else
    unset DEBUG
fi

if [[ ${FLAG_BUILD_RUN} -eq 1 ]]; then
    exec_build_run
elif [[ ${FLAG_TESTS} -eq 1 ]]; then
    exec_tests
fi

exit 0
