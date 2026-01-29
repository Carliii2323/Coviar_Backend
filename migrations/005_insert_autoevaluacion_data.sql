-- Datos de Prueba para el Sistema de Autoevaluación

-- Insertar segmentos
INSERT INTO segmentos (nombre, min_turistas, max_turistas) VALUES
  ('Pequeña bodega', 100, 500),
  ('Bodega mediana', 500, 2000),
  ('Bodega grande', 2000, NULL);

-- Insertar capítulos
INSERT INTO capitulos (nombre, descripcion, orden) VALUES
  ('Contexto de la organización', 'Públicos de interés o partes interesadas relevantes a las actividades enoturísticas', 1),
  ('Planificación', 'Planificación estratégica y operativa', 2),
  ('Implementación', 'Implementación de políticas y programas', 3),
  ('Seguimiento y mejora', 'Monitoreo y mejora continua', 4);

-- Insertar indicadores
INSERT INTO indicadores (id_capitulo, nombre, descripcion, orden) VALUES
  (1, 'Indicador 1.1', 'Tendencias del turismo vitivinícola', 1),
  (1, 'Indicador 1.2', 'Identificación de públicos de interés', 2),
  (2, 'Indicador 2.1', 'Existencia de plan estratégico', 1),
  (2, 'Indicador 2.2', 'Metas cuantificables', 2),
  (3, 'Indicador 3.1', 'Implementación de políticas', 1),
  (3, 'Indicador 3.2', 'Documentación de procesos', 2),
  (4, 'Indicador 4.1', 'Sistema de indicadores', 1),
  (4, 'Indicador 4.2', 'Mejora continua', 2);

-- Insertar niveles de respuesta
INSERT INTO niveles_respuesta (id_indicador, nombre, descripcion, puntos) VALUES
  (1, 'No alcanza nivel 1', 'Sin conocimiento o evidencia', 0),
  (1, 'Nivel 1', 'Conocimiento inicial', 1),
  (1, 'Nivel 2', 'Conocimiento intermedio', 2),
  (1, 'Nivel 3', 'Conocimiento avanzado', 3),
  (1, 'Nivel 4', 'Liderazgo en la materia', 4),
  
  (2, 'No alcanza nivel 1', 'Sin conocimiento o evidencia', 0),
  (2, 'Nivel 1', 'Conocimiento inicial', 1),
  (2, 'Nivel 2', 'Conocimiento intermedio', 2),
  (2, 'Nivel 3', 'Conocimiento avanzado', 3),
  (2, 'Nivel 4', 'Liderazgo en la materia', 4),
  
  (3, 'No alcanza nivel 1', 'Sin plan', 0),
  (3, 'Nivel 1', 'Plan básico', 1),
  (3, 'Nivel 2', 'Plan estructurado', 2),
  (3, 'Nivel 3', 'Plan alineado con sostenibilidad', 3),
  (3, 'Nivel 4', 'Plan de excelencia', 4),
  
  (4, 'No alcanza nivel 1', 'Sin metas', 0),
  (4, 'Nivel 1', 'Metas genéricas', 1),
  (4, 'Nivel 2', 'Metas cuantificables', 2),
  (4, 'Nivel 3', 'Metas con indicadores', 3),
  (4, 'Nivel 4', 'Metas alineadas con ODS', 4),
  
  (5, 'No alcanza nivel 1', 'Sin implementación', 0),
  (5, 'Nivel 1', 'Implementación parcial', 1),
  (5, 'Nivel 2', 'Implementación en proceso', 2),
  (5, 'Nivel 3', 'Implementación completa', 3),
  (5, 'Nivel 4', 'Implementación con mejoras', 4),
  
  (6, 'No alcanza nivel 1', 'Sin documentación', 0),
  (6, 'Nivel 1', 'Documentación básica', 1),
  (6, 'Nivel 2', 'Documentación estructurada', 2),
  (6, 'Nivel 3', 'Documentación completa', 3),
  (6, 'Nivel 4', 'Documentación actualizada', 4),
  
  (7, 'No alcanza nivel 1', 'Sin indicadores', 0),
  (7, 'Nivel 1', 'Indicadores básicos', 1),
  (7, 'Nivel 2', 'Sistema de indicadores', 2),
  (7, 'Nivel 3', 'Sistema integrado', 3),
  (7, 'Nivel 4', 'Sistema de excelencia', 4),
  
  (8, 'No alcanza nivel 1', 'Sin mejora', 0),
  (8, 'Nivel 1', 'Mejora esporádica', 1),
  (8, 'Nivel 2', 'Mejora continua', 2),
  (8, 'Nivel 3', 'Mejora sistemática', 3),
  (8, 'Nivel 4', 'Mejora estratégica', 4);

-- Asociar indicadores con segmentos
-- Pequeña bodega (id_segmento = 1): indicadores 1, 2, 3, 4
INSERT INTO segmento_indicador (id_segmento, id_indicador) VALUES
  (1, 1), (1, 2), (1, 3), (1, 4);

-- Bodega mediana (id_segmento = 2): todos los indicadores
INSERT INTO segmento_indicador (id_segmento, id_indicador) VALUES
  (2, 1), (2, 2), (2, 3), (2, 4), (2, 5), (2, 6), (2, 7), (2, 8);

-- Bodega grande (id_segmento = 3): todos los indicadores
INSERT INTO segmento_indicador (id_segmento, id_indicador) VALUES
  (3, 1), (3, 2), (3, 3), (3, 4), (3, 5), (3, 6), (3, 7), (3, 8);

-- Insertar niveles de sostenibilidad
INSERT INTO niveles_sostenibilidad (id_segmento, nombre, min_puntaje, max_puntaje) VALUES
  (1, 'En desarrollo', 0, 7),
  (1, 'Emergente', 8, 14),
  (1, 'Consolidado', 15, 20),
  
  (2, 'En desarrollo', 0, 15),
  (2, 'Emergente', 16, 31),
  (2, 'Consolidado', 32, 50),
  
  (3, 'En desarrollo', 0, 21),
  (3, 'Emergente', 22, 43),
  (3, 'Consolidado', 44, 64);

-- Nota: Ejecutar estos comandos en Supabase SQL Editor después de ejecutar 
-- la migración 004_create_autoevaluacion_tables.sql
