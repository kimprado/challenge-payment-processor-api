#!/bin/bash

set -e

CLI_MODE="-n"
GUI_MODE=""
MODE=""

PATH_TEST_PLAN="./jmeter-load-test-plan.jmx"
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

if  [ ! -e $PATH_TEST_PLAN ]; then
    PATH_TEST_PLAN="./test/jmeter-load-test-plan.jmx"
fi

if  [ ! -e $TEST_CASES_PATH ]; then
    TEST_CASES_PATH="./test/jmeter-requests-config.csv"
fi

export JMETER_API_HOST=localhost
export JMETER_API_PORT=80
export JMETER_API_CONTEXT=/api
export JMETER_TEST_CASES_PATH=$TEST_CASES_PATH

export HEAP="-Xms3g -Xmx4g -XX:MaxMetaspaceSize=3g -Xmn2g"

echo "Aguardando Target..."
# Esperar startup
echo "'Aguardo' encerrado"

jmeter $MODE -t $PATH_TEST_PLAN

# Verificar status do encerramento do Teste Jmeter, com "$?"
if [ $? -eq 0 ]
then
    exit 0
else    
    exit 1
fi
