#!/bin/bash

set -e

CLI_MODE="-n"
GUI_MODE=""
MODE=""

TEST_PLAN_PATH="./jmeter-load-test-plan.jmx"
TEST_CASES_PATH="./jmeter-requests-config.csv"

if [ -z "$1" ]
then
    echo "JMeter in CLI mode"
    MODE=$CLI_MODE
    export JMETER_EXIT_ON_ERROR=true
else
    echo "JMeter in GUI mode"
    MODE=$GUI_MODE
    export JMETER_EXIT_ON_ERROR=false
fi

if  [ ! -e $TEST_PLAN_PATH ]; then
    TEST_PLAN_PATH="./test/jmeter-load-test-plan.jmx"
fi

if  [ ! -e $TEST_CASES_PATH ]; then
    TEST_CASES_PATH="./test/jmeter-requests-config.csv"
fi

export JMETER_API_HOST=localhost
export JMETER_API_PORT=80
export JMETER_API_CONTEXT=/api
export JMETER_TEST_CASES_PATH=$TEST_CASES_PATH

export HEAP="-Xms3g -Xmx4g -XX:MaxMetaspaceSize=3g -Xmn2g"

jmeter $MODE -t $TEST_PLAN_PATH

# Verificar status do encerramento do Teste Jmeter, com "$?"
if [ $? -eq 0 ]
then
    exit 0
else    
    exit 1
fi
