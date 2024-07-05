#!/usr/bin/env bash

# GNU's nm is required in order for exercises where binary stripping
# to be testable in a robust way.
check_nm() {
    if ! command -v nm > /dev/null 2>&1; then
        echo "STARTUP: error: nm is not installed"
        exit 1
    fi

    if ! nm --version 2>&1 | grep -q "GNU"; then 
        echo "STARTUP: error: GNU nm is required" >&2
        exit 1
    fi
    echo "STARTUP: GNU nm is installed"
}

echo "STARTUP: all dependencies are satisfied"
check_nm
