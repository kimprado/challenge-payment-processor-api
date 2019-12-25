#!/bin/bash

set -e

echo "Executando teste de carga"

ab -n 61000 -c 103 -p POST_Transaction.json -T application/json -H 'X-ACQUIRER-ID: Stone' -m POST http://localhost:80/api/transactions

# Verificar status do encerramento do teste, com "$?"
if [ $? -eq 0 ]
then
    exit 0
else    
    exit 1
fi
