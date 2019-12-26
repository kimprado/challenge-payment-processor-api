#!/bin/bash

set -e

echo "Executando teste de carga"

REQUEST_BODY_FILE="./apachebench-POST_Transaction.json"
if  [ ! -e $REQUEST_BODY_FILE ]; then
    REQUEST_BODY_FILE="./test/apachebench-POST_Transaction.json"
fi

ab -n 610000 -c 103 -p $REQUEST_BODY_FILE -T application/json -H 'X-ACQUIRER-ID: Stone' -m POST http://localhost:80/api/transactions

# Verificar status do encerramento do teste, com "$?"
if [ $? -eq 0 ]
then
    exit 0
else    
    exit 1
fi
