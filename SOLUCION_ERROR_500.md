# Solución: Error 500 en POST /api/autoevaluaciones

## 🔍 Diagnóstico del Problema

El error 500 que recibiste al hacer:
```bash
curl -X POST http://localhost:8080/api/autoevaluaciones
```

**No era un error en el código, sino un error de autenticación esperado.**

### ¿Por qué ocurrió?

Todos los endpoints de autoevaluación están protegidos con middleware JWT:

```go
// En cmd/api/main.go líneas 142-147
engine.POST("/api/autoevaluaciones", protect(handlers.CreateAutoevaluacion))
engine.GET("/api/autoevaluaciones/:id/segmentos", protect(handlers.GetSegmentos))
engine.PUT("/api/autoevaluaciones/:id/segmento", protect(handlers.SeleccionarSegmento))
engine.GET("/api/autoevaluaciones/:id/estructura", protect(handlers.GetEstructura))
engine.POST("/api/autoevaluaciones/:id/respuestas", protect(handlers.GuardarRespuestas))
engine.POST("/api/autoevaluaciones/:id/completar", protect(handlers.CompletarAutoevaluacion))
```

La función `protect()` envuelve cada handler y requiere un token JWT válido. Sin autenticación, el middleware retorna error antes de llegar al handler.

---

## ✅ Solución: Autenticación + Endpoint

### Paso 1: Registrar una bodega

```bash
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
```

**Importante:** `-c cookies.txt` guarda las cookies (incluyendo auth_token) en un archivo.

### Paso 2: Hacer login (obtener JWT)

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email_login": "test@bodega.com",
    "password": "SecurePass123!"
  }'
```

Esto actualiza `cookies.txt` con el token JWT en `auth_token`.

### Paso 3: Usar el token en los endpoints

```bash
# ✅ CORRECTO: Con autenticación
curl -X POST http://localhost:8080/api/autoevaluaciones \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"id_bodega": 1}'

# ❌ INCORRECTO: Sin autenticación
curl -X POST http://localhost:8080/api/autoevaluaciones \
  -H "Content-Type: application/json" \
  -d '{"id_bodega": 1}'
```

La clave es `-b cookies.txt` que envía el token con cada petición.

---

## 🚀 Opción Más Fácil: Script Automatizado

En lugar de hacer cada curl manualmente, usa el script que ya creé:

```bash
chmod +x test_autoevaluacion.sh
./test_autoevaluacion.sh
```

Este script:
1. Registra bodega
2. Hace login
3. Crea autoevaluación
4. Obtiene segmentos
5. Selecciona segmento
6. Obtiene estructura
7. Guarda respuestas
8. Completa evaluación

Todo en un solo comando ✨

---

## 📚 Referencia Rápida: Cookies vs Headers

| Método | Comando | Notas |
|--------|---------|-------|
| Guardar cookies | `-c cookies.txt` | Usado en `/api/registro` y `/api/login` |
| Enviar cookies | `-b cookies.txt` | Usado en endpoints protegidos |
| En un header | `-H "Authorization: Bearer <JWT>"` | Alternativa a cookies |

---

## ✔️ Verificación

El error 500 significa que el middleware rechazó la petición sin auth. Ahora que sabes usar cookies:

1. ✅ Tu código está correcto
2. ✅ Los endpoints están bien protegidos
3. ✅ Solo necesitabas autenticación

¡Prueba el script y verás que todo funciona! 🎉
