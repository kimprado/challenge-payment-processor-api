#!/bin/bash

set -e
clear

CLI_MODE="-n"
GUI_MODE=""
MODE=""

PATH_TEST_PLAN="./jmeter-load-test-plan.jmx"

if [ -z "$1" ]
then
    echo "JMeter in CLI mode"
    MODE=$CLI_MODE
else
    echo "JMeter in GUI mode"
    MODE=$GUI_MODE
fi

if  [ ! -e $PATH_TEST_PLAN ]; then
    PATH_TEST_PLAN="./test/jmeter-load-test-plan.jmx"
fi

export JMETER_API_HOST=localhost
export JMETER_API_PORT=80
export JMETER_API_CONTEXT=/api

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
