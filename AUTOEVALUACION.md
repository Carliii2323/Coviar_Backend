# Sistema de Autoevaluación - Documentación

## Descripción General

El sistema de autoevaluación permite que las cuentas de tipo **BODEGA** completen un cuestionario estructurado sobre sostenibilidad. El proceso consta de 6 pasos principales.

---

## Flujo de Autoevaluación

### 1. Crear Autoevaluación

**Método:** `POST`  
**Endpoint:** `/api/autoevaluaciones`  
**Requiere:** Autenticación JWT (cookie `auth_token`)

**Importante:** Debes estar autenticado (hacer login primero) para acceder a este endpoint.

**Body:**
```json
{
  "id_bodega": 1
}
```

**Respuesta (201 Created):**
```json
{
  "id_autoevaluacion": 42,
  "fecha_inicio": "2026-01-27T10:30:00Z",
  "fecha_fin": null,
  "estado": "PENDIENTE",
  "id_bodega": 1,
  "id_segmento": null
}
```

---

### 2. Obtener Segmentos Disponibles

**Método:** `GET`  
**Endpoint:** `/api/autoevaluaciones/{id_autoevaluacion}/segmentos`  
**Requiere:** Autenticación

**Respuesta (200 OK):**
```json
[
  {
    "id_segmento": 1,
    "nombre": "Pequeña bodega",
    "min_turistas": 100,
    "max_turistas": 500
  },
  {
    "id_segmento": 2,
    "nombre": "Bodega mediana",
    "min_turistas": 500,
    "max_turistas": 2000
  }
]
```

---

### 3. Seleccionar Segmento

**Método:** `PUT`  
**Endpoint:** `/api/autoevaluaciones/{id_autoevaluacion}/segmento`  
**Requiere:** Autenticación

**Body:**
```json
{
  "id_segmento": 2
}
```

**Respuesta (200 OK):**
```json
{
  "mensaje": "Segmento seleccionado correctamente"
}
```

---

### 4. Obtener Estructura del Cuestionario

**Método:** `GET`  
**Endpoint:** `/api/autoevaluaciones/{id_autoevaluacion}/estructura`  
**Requiere:** Autenticación

**Respuesta (200 OK):**
```json
{
  "capitulos": [
    {
      "capitulo": {
        "id_capitulo": 1,
        "nombre": "Capítulo 1",
        "descripcion": "Contexto de la organización",
        "orden": 1
      },
      "indicadores": [
        {
          "indicador": {
            "id_indicador": 1,
            "id_capitulo": 1,
            "nombre": "Indicador 1.1",
            "descripcion": "Tendencias del turismo vitivinícola",
            "orden": 1
          },
          "niveles_respuesta": [
            {
              "id_nivel_respuesta": 9,
              "id_indicador": 1,
              "nombre": "No alcanza nivel 1",
              "descripcion": "Descripción del nivel",
              "puntos": 0
            },
            {
              "id_nivel_respuesta": 10,
              "id_indicador": 1,
              "nombre": "Nivel 1",
              "descripcion": "Descripción del nivel",
              "puntos": 1
            }
          ],
          "habilitado": true
        }
      ]
    }
  ]
}
```

**Nota sobre "habilitado":** Este campo es `true` si el indicador está asociado al segmento seleccionado, `false` en caso contrario.

---

### 5. Guardar Respuestas

**Método:** `POST`  
**Endpoint:** `/api/autoevaluaciones/{id_autoevaluacion}/respuestas`  
**Requiere:** Autenticación

**Body:**
```json
{
  "respuestas": [
    {
      "id_indicador": 1,
      "id_nivel_respuesta": 10
    },
    {
      "id_indicador": 2,
      "id_nivel_respuesta": 15
    }
  ]
}
```

**Respuesta (200 OK):**
```json
{
  "mensaje": "Respuestas guardadas correctamente"
}
```

**Notas:**
- Las respuestas anteriores se limpian automáticamente
- Cada respuesta vincula un indicador con un nivel de respuesta
- Los puntos se calculan a partir del nivel_respuesta seleccionado

---

### 6. Completar Autoevaluación

**Método:** `POST`  
**Endpoint:** `/api/autoevaluaciones/{id_autoevaluacion}/completar`  
**Requiere:** Autenticación

**Respuesta (200 OK):**
```json
{
  "mensaje": "Autoevaluación completada correctamente"
}
```

**Validaciones:**
- La autoevaluación debe tener al menos una respuesta guardada
- Se establece el estado a "COMPLETADA"
- Se registra la fecha_fin

---

## Modelos de Base de Datos

### Segmentos
```sql
CREATE TABLE segmentos (
    id_segmento INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nombre TEXT NOT NULL,
    min_turistas INTEGER NOT NULL,
    max_turistas INTEGER
);
```

### Autoevaluaciones
```sql
CREATE TABLE autoevaluaciones (
    id_autoevaluacion INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    fecha_inicio TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fecha_fin TIMESTAMPTZ,
    estado ENUM ('PENDIENTE', 'COMPLETADA', 'CANCELADA') DEFAULT 'PENDIENTE',
    id_bodega INTEGER NOT NULL REFERENCES bodegas(id_bodega),
    id_segmento INTEGER REFERENCES segmentos(id_segmento)
);
```

### Capítulos
```sql
CREATE TABLE capitulos (
    id_capitulo INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nombre TEXT NOT NULL,
    descripcion TEXT NOT NULL,
    orden INTEGER NOT NULL
);
```

### Indicadores
```sql
CREATE TABLE indicadores (
    id_indicador INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    id_capitulo INTEGER NOT NULL REFERENCES capitulos(id_capitulo),
    nombre TEXT NOT NULL,
    descripcion TEXT NOT NULL,
    orden INTEGER NOT NULL
);
```

### Niveles de Respuesta
```sql
CREATE TABLE niveles_respuesta (
    id_nivel_respuesta INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    id_indicador INTEGER NOT NULL REFERENCES indicadores(id_indicador),
    nombre TEXT NOT NULL,
    descripcion TEXT NOT NULL,
    puntos INTEGER NOT NULL
);
```

### Respuestas
```sql
CREATE TABLE respuestas (
    id_respuesta INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    id_nivel_respuesta INTEGER NOT NULL REFERENCES niveles_respuesta(id_nivel_respuesta),
    id_indicador INTEGER NOT NULL REFERENCES indicadores(id_indicador),
    id_autoevaluacion INTEGER NOT NULL REFERENCES autoevaluaciones(id_autoevaluacion)
);
```

### Segmento-Indicador (Relación)
```sql
CREATE TABLE segmento_indicador (
    id_segmento INTEGER NOT NULL REFERENCES segmentos(id_segmento),
    id_indicador INTEGER NOT NULL REFERENCES indicadores(id_indicador),
    PRIMARY KEY (id_segmento, id_indicador)
);
```

---

## Arquitectura

El sistema sigue la arquitectura de capas del proyecto:

- **Domain**: Modelos (`domain/models.go`)
- **Repository**: Acceso a datos (`internal/repository/postgres/`)
  - `segmento_repository.go`
  - `autoevaluacion_repository.go`
  - `capitulo_repository.go`
  - `indicador_repository.go`
  - `nivel_respuesta_repository.go`
  - `respuesta_repository.go`
- **Service**: Lógica de negocio (`internal/service/autoevaluacion_service.go`)
- **Handler**: Controladores HTTP (`internal/handler/autoevaluacion_handler.go`)
- **Router**: Definición de rutas (`cmd/api/main.go`)

---

## Ejemplo Completo de Flujo

```bash
# 1. Crear autoevaluación
curl -X POST http://localhost:8080/api/autoevaluaciones \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=<JWT_TOKEN>" \
  -d '{"id_bodega": 1}'

# Respuesta:
# {"id_autoevaluacion": 42, "estado": "PENDIENTE", ...}

# 2. Obtener segmentos
curl -X GET http://localhost:8080/api/autoevaluaciones/42/segmentos \
  -H "Cookie: auth_token=<JWT_TOKEN>"

# 3. Seleccionar segmento
curl -X PUT http://localhost:8080/api/autoevaluaciones/42/segmento \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=<JWT_TOKEN>" \
  -d '{"id_segmento": 2}'

# 4. Obtener estructura
curl -X GET http://localhost:8080/api/autoevaluaciones/42/estructura \
  -H "Cookie: auth_token=<JWT_TOKEN>"

# 5. Guardar respuestas
curl -X POST http://localhost:8080/api/autoevaluaciones/42/respuestas \
  -H "Content-Type: application/json" \
  -H "Cookie: auth_token=<JWT_TOKEN>" \
  -d '{"respuestas": [{"id_indicador": 1, "id_nivel_respuesta": 10}]}'

# 6. Completar autoevaluación
curl -X POST http://localhost:8080/api/autoevaluaciones/42/completar \
  -H "Cookie: auth_token=<JWT_TOKEN>"
```

---

## Consideraciones de Seguridad

- Todos los endpoints requieren autenticación JWT
- Las cookies son HttpOnly y SameSite Lax
- Solo las bodegas autenticadas pueden crear autoevaluaciones
- Se valida que los indicadores pertenezcan al segmento seleccionado

---

## Estados de Autoevaluación

- **PENDIENTE**: Autoevaluación en progreso
- **COMPLETADA**: Autoevaluación finalizada
- **CANCELADA**: Autoevaluación cancelada (no implementado aún)
