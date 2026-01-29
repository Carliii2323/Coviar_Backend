# Sistema de Autoevaluación - Guía de Instalación y Uso

## 📦 Requisitos

- Go 1.19+
- PostgreSQL (vía Supabase)
- Variables de entorno configuradas en `.env`

---

## 🔧 Pasos de Instalación

### 1. Ejecutar Migraciones SQL

Conectarse a Supabase y ejecutar en el SQL Editor:

```bash
# Primero ejecutar la migración de tablas
# Archivo: migrations/004_create_autoevaluacion_tables.sql

# Luego ejecutar los datos de prueba
# Archivo: migrations/005_insert_autoevaluacion_data.sql
```

O desde la terminal (si tienes acceso a psql):
```bash
psql -h db.supabase.co -U postgres -d postgres \
  -f migrations/004_create_autoevaluacion_tables.sql \
  -f migrations/005_insert_autoevaluacion_data.sql
```

### 2. Verificar Compilación

```bash
cd /path/to/Coviar_Backend
go build -o ./bin/api ./cmd/api
```

Si compila sin errores, está listo.

### 3. Ejecutar el Servidor

```bash
# Opción 1: Compilado
./bin/api

# Opción 2: Directo con go run
go run ./cmd/api
```

Debería ver:
```
✓ Configuración cargada
✓ Conexión a Supabase establecida
✓ Repositorios inicializados
✓ Servicios inicializados
✓ Handlers inicializados
🚀 Servidor iniciando en http://0.0.0.0:8080
```

---

## 🧪 Pruebas Rápidas

### ⚠️ IMPORTANTE: Autenticación Requerida

**Todos los endpoints de autoevaluación requieren autenticación JWT.** Debes completar el flujo de login antes de acceder a estos endpoints.

### Opción 1: Script Automatizado (RECOMENDADO)

```bash
chmod +x test_autoevaluacion.sh
./test_autoevaluacion.sh
```

Este script automáticamente:
1. ✅ Registra una nueva bodega
2. ✅ Hace login y obtiene token JWT
3. ✅ Crea una autoevaluación
4. ✅ Obtiene segmentos disponibles
5. ✅ Selecciona un segmento
6. ✅ Obtiene la estructura del cuestionario

### Opción 2: Usando curl manualmente

Importante: **Guardar cookies con `-c cookies.txt` en login y enviarlas con `-b cookies.txt` en cada petición**

```bash
# 1. Registrar una bodega
curl -X POST http://localhost:8080/api/registro \
  -H "Content-Type: application/json" \
  -c cookies.txt \
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
  }'

# 2. Hacer login (guarda token en cookies.txt)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email_login": "test@bodega.com",
    "password": "SecurePass123!"
  }'

# 3. Crear autoevaluación (usar -b para enviar cookies)
AUTOEVALID=$(curl -s -X POST http://localhost:8080/api/autoevaluaciones \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"id_bodega": 1}' | jq -r '.id_autoevaluacion')

echo "✅ Autoevaluación creada: $AUTOEVALID"

# 4. Obtener segmentos
curl -X GET http://localhost:8080/api/autoevaluaciones/$AUTOEVALID/segmentos \
  -H "Content-Type: application/json" \
  -b cookies.txt

# 5. Seleccionar segmento
curl -X PUT http://localhost:8080/api/autoevaluaciones/$AUTOEVALID/segmento \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"id_segmento": 2}'

# 6. Obtener estructura
curl -X GET http://localhost:8080/api/autoevaluaciones/$AUTOEVALID/estructura \
  -H "Content-Type: application/json" \
  -b cookies.txt

# 7. Guardar respuestas
curl -X POST http://localhost:8080/api/autoevaluaciones/$AUTOEVALID/respuestas \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "respuestas": [
      {"id_indicador": 1, "id_nivel_respuesta": 2},
      {"id_indicador": 2, "id_nivel_respuesta": 3}
    ]
  }'

# 8. Completar autoevaluación
curl -X POST http://localhost:8080/api/autoevaluaciones/$AUTOEVALID/completar \
  -H "Content-Type: application/json" \
  -b cookies.txt
```

**Clave:** La opción `-b cookies.txt` envía las cookies de autenticación con cada petición. Sin esto, recibirás error 401 o 500.

---

## 📊 Estructura de Datos

### Relaciones Principales

```
Bodega (id_bodega)
  └─ Autoevaluación (id_autoevaluacion, id_bodega)
       └─ Segmento (id_segmento)
            └─ Indicadores habilitados (segmento_indicador)
                 └─ Respuestas (id_respuesta)
                      └─ Nivel de Respuesta (id_nivel_respuesta, puntos)
```

### Datos de Prueba

- **3 Segmentos** (Pequeña, Mediana, Grande)
- **4 Capítulos** (Contexto, Planificación, Implementación, Seguimiento)
- **8 Indicadores** (2 por capítulo)
- **5 Niveles de Respuesta** por indicador (0-4 puntos)

---

## 🔐 Autenticación

### Obtener Token JWT

```bash
# Primero, registrar una bodega
curl -X POST http://localhost:8080/api/registro \
  -H "Content-Type: application/json" \
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
  }'

# Luego, hacer login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email_login": "test@bodega.com",
    "password": "SecurePass123!"
  }'

# El token se guarda en cookies.txt (HttpOnly)
# Usarlo en todas las requests de autoevaluación:
# -H "Cookie: auth_token=<token>"
```

---

## 🐛 Troubleshooting

### Error: "No se conecta a Supabase"
```
Verificar:
- SUPABASE_URL en .env
- SUPABASE_KEY en .env
- SUPABASE_DB_PASSWORD en .env
- Conexión a internet
```

### Error: "Tabla no existe"
```
Ejecutar migraciones:
1. migrations/004_create_autoevaluacion_tables.sql
2. migrations/005_insert_autoevaluacion_data.sql
```

### Error: "Unauthorized"
```
- El JWT token puede estar expirado
- Hacer login nuevamente
- Verificar que el token esté en las cookies
```

### Error: "Segmento no encontrado"
```
- Verificar que id_segmento existe en la BD
- Usar GET /api/autoevaluaciones/{id}/segmentos para listar
```

---

## 📈 Próximas Mejoras

1. **Cálculo de Puntaje**: Sumar puntos de respuestas y clasificar en nivel de sostenibilidad
2. **Reportes**: Generar reportes PDF con resultados
3. **Historial**: Permitir múltiples autoevaluaciones por bodega
4. **Comparativa**: Comparar resultados entre segmentos
5. **Recomendaciones**: Sugerir mejoras basadas en respuestas

---

## 📚 Documentación Adicional

- [AUTOEVALUACION.md](./AUTOEVALUACION.md) - Especificación de endpoints
- [IMPLEMENTACION_AUTOEVALUACION.md](./IMPLEMENTACION_AUTOEVALUACION.md) - Detalles técnicos

---

## ✅ Checklist de Implementación

- [x] Modelos de dominio creados
- [x] Repositorios PostgreSQL implementados
- [x] Servicio de autoevaluación desarrollado
- [x] Handler HTTP creado
- [x] Rutas registradas en router
- [x] Migraciones SQL generadas
- [x] Datos de prueba insertados
- [x] Compilación exitosa
- [x] Documentación completa

---

## 📞 Soporte

Para preguntas sobre el sistema de autoevaluación:
1. Revisar `AUTOEVALUACION.md`
2. Revisar `IMPLEMENTACION_AUTOEVALUACION.md`
3. Revisar ejemplos curl en esta guía
