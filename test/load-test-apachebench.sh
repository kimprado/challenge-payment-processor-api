#!/bin/bash

set -e
clear

echo "Executando teste de carga"

# Verificar status do encerramento do teste, com "$?"
if [ $? -eq 0 ]
then
    exit 0
else    
    exit 1
fi
