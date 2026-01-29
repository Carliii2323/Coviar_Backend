-- Migración: Crear tablas de autoevaluación
-- Fecha: 2026-01-27

-- Crear enum para estado de autoevaluación
CREATE TYPE estado_autoevaluacion AS ENUM (
  'PENDIENTE', 
  'COMPLETADA', 
  'CANCELADA'
);

-- Tabla de segmentos
CREATE TABLE IF NOT EXISTS segmentos (
    id_segmento integer generated always as identity,
    nombre text not null,
    min_turistas integer not null,
    max_turistas integer,
    constraint segmentos_pkey primary key (id_segmento)
);

-- Tabla de niveles de sostenibilidad
CREATE TABLE IF NOT EXISTS niveles_sostenibilidad (
    id_nivel_sostenibilidad integer generated always as identity,
    id_segmento integer not null,
    nombre text not null,
    min_puntaje integer not null,
    max_puntaje integer not null,
    constraint niveles_sostenibilidad_pkey primary key (id_nivel_sostenibilidad),
    constraint niveles_sostenibilidad_segmento_fk foreign key (id_segmento) references segmentos (id_segmento)
);

-- Tabla de capítulos
CREATE TABLE IF NOT EXISTS capitulos (
    id_capitulo integer generated always as identity,
    nombre text not null,
    descripcion text not null,
    orden integer not null,
    constraint capitulos_pkey primary key (id_capitulo)
);

-- Tabla de indicadores
CREATE TABLE IF NOT EXISTS indicadores (
    id_indicador integer generated always as identity,
    id_capitulo integer not null,
    nombre text not null,
    descripcion text not null,
    orden integer not null,
    constraint indicadores_pkey primary key (id_indicador),
    constraint indicadores_capitulo_fk foreign key (id_capitulo) references capitulos (id_capitulo)
);

-- Tabla de niveles de respuesta
CREATE TABLE IF NOT EXISTS niveles_respuesta (
    id_nivel_respuesta integer generated always as identity,
    id_indicador integer not null,
    nombre text not null,
    descripcion text not null,
    puntos integer not null,
    constraint niveles_respuesta_pk primary key (id_nivel_respuesta),
    constraint niveles_respuesta_indicador_fk foreign key (id_indicador) references indicadores (id_indicador)
);

-- Tabla de autoevaluaciones
CREATE TABLE IF NOT EXISTS autoevaluaciones (
    id_autoevaluacion integer generated always as identity,
    fecha_inicio timestamptz not null default now(),
    fecha_fin timestamptz,
    estado estado_autoevaluacion not null default 'PENDIENTE',
    id_bodega integer not null,
    id_segmento integer,
    constraint autoevaluaciones_pk primary key (id_autoevaluacion),
    constraint autoevaluaciones_bodega_fk foreign key (id_bodega) references bodegas (id_bodega),
    constraint autoevaluaciones_segmento_fk foreign key (id_segmento) references segmentos (id_segmento)
);

-- Tabla de respuestas
CREATE TABLE IF NOT EXISTS respuestas (
    id_respuesta integer generated always as identity,
    id_nivel_respuesta integer not null,
    id_indicador integer not null,
    id_autoevaluacion integer not null,
    constraint respuestas_pk primary key (id_respuesta),
    constraint respuestas_nivel_respuesta_fk foreign key (id_nivel_respuesta) references niveles_respuesta (id_nivel_respuesta),
    constraint respuestas_indicador_fk foreign key (id_indicador) references indicadores (id_indicador),
    constraint respuestas_autoevaluacion_fk foreign key (id_autoevaluacion) references autoevaluaciones (id_autoevaluacion)
);

-- Tabla de relación segmento-indicador
CREATE TABLE IF NOT EXISTS segmento_indicador (
    id_segmento integer not null,
    id_indicador integer not null,
    constraint segmento_indicador_pk primary key(id_segmento, id_indicador),
    constraint segmento_indicador_segmento_fk foreign key (id_segmento) references segmentos (id_segmento),
    constraint segmento_indicador_indicador_fk foreign key (id_indicador) references indicadores (id_indicador)
);

-- Crear índices para mejor rendimiento
CREATE INDEX IF NOT EXISTS idx_autoevaluaciones_bodega ON autoevaluaciones(id_bodega);
CREATE INDEX IF NOT EXISTS idx_autoevaluaciones_segmento ON autoevaluaciones(id_segmento);
CREATE INDEX IF NOT EXISTS idx_autoevaluaciones_estado ON autoevaluaciones(estado);
CREATE INDEX IF NOT EXISTS idx_respuestas_autoevaluacion ON respuestas(id_autoevaluacion);
CREATE INDEX IF NOT EXISTS idx_indicadores_capitulo ON indicadores(id_capitulo);
CREATE INDEX IF NOT EXISTS idx_niveles_respuesta_indicador ON niveles_respuesta(id_indicador);
