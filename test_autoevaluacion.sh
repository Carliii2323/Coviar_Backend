#!/bin/bash

# Script de prueba para autoevaluaciones

API="http://localhost:8080"
COOKIES="/tmp/cookies.txt"

echo "🔐 Paso 1: Registrar bodega..."
REGISTRO=$(curl -s -X POST $API/api/registro \
  -H "Content-Type: application/json" \
  -c "$COOKIES" \
  -d '{
    "bodega": {
      "razon_social": "Bodega Test",
      "nombre_fantasia": "Test",
      "cuit": "12345678901",
      "calle": "Calle Principal",
      "numeracion": "123",
      "id_localidad": 1,
      "telefono": "1234567890",
      "email_institucional": "test@bodega.com"
    },
    "cuenta": {
      "email_login": "test@bodega.com",
      "password": "SecurePass123!"
    },
    "responsable": {
      "nombre": "Juan",
      "apellido": "Pérez",
      "cargo": "Gerente",
      "dni": "12345678"
    }
  }')

echo "Respuesta: $REGISTRO"
ID_BODEGA=$(echo "$REGISTRO" | jq -r '.id_bodega')
echo "ID Bodega: $ID_BODEGA"

echo ""
echo "🔐 Paso 2: Login..."
LOGIN=$(curl -s -X POST $API/api/login \
  -H "Content-Type: application/json" \
  -c "$COOKIES" \
  -d '{
    "email_login": "test@bodega.com",
    "password": "SecurePass123!"
  }')

echo "Respuesta: $LOGIN"

echo ""
echo "✅ Paso 3: Crear autoevaluación..."
AUTO=$(curl -s -X POST $API/api/autoevaluaciones \
  -H "Content-Type: application/json" \
  -b "$COOKIES" \
  -d "{\"id_bodega\": $ID_BODEGA}")

echo "Respuesta: $AUTO"
ID_AUTO=$(echo "$AUTO" | jq -r '.id_autoevaluacion')
echo "ID Autoevaluación: $ID_AUTO"

echo ""
echo "✅ Paso 4: Obtener segmentos..."
curl -s -X GET $API/api/autoevaluaciones/$ID_AUTO/segmentos \
  -H "Content-Type: application/json" \
  -b "$COOKIES" | jq .

echo ""
echo "✅ Paso 5: Seleccionar segmento..."
curl -s -X PUT $API/api/autoevaluaciones/$ID_AUTO/segmento \
  -H "Content-Type: application/json" \
  -b "$COOKIES" \
  -d '{"id_segmento": 2}' | jq .

echo ""
echo "✅ Paso 6: Obtener estructura..."
curl -s -X GET $API/api/autoevaluaciones/$ID_AUTO/estructura \
  -H "Content-Type: application/json" \
  -b "$COOKIES" | jq . | head -50

echo ""
echo "✅ Script completado"
