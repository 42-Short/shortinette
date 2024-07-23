#!/usr/bin/env sh

check_nm() {
    if ! command -v nm > /dev/null 2>&1; then
        echo "STARTUP ERROR: nm is not installed"
        exit 1
    fi

    if ! nm --version 2>&1 | grep -q "GNU"; then 
        echo "STARTUP ERROR: GNU nm is required" >&2
        exit 1
    fi
    echo "STARTUP: GNU nm is installed"
}

build_testenv() {
    if ! docker build -t testenv .; then 
        echo "STARTUP ERROR: could not build testenv from ./Dockerfile"
        exit 1
    fi
}

check_testenv() {
    if ! docker image ls | grep testenv > /dev/null 2>&1; then
        echo "STARTUP ERROR: 'testenv' Docker image missing"
        echo "STARTUP: trying to build from ./Dockerfile"
        build_testenv
    fi
    echo "STARTUP: 'testenv' Docker image found"
}


check_nm
check_testenv
echo "STARTUP: all dependencies are satisfied"
