-- MIGRACIÓN: Cambios de estructura para bodegas, cuentas y responsables según nueva especificación

-- Cambios en tabla bodegas
ALTER TABLE bodega RENAME TO bodegas;
ALTER TABLE bodegas ALTER COLUMN cuit TYPE CHAR(11) USING cuit::char(11);
ALTER TABLE bodegas ALTER COLUMN inv_bod TYPE CHAR(6) USING inv_bod::char(6);
ALTER TABLE bodegas ALTER COLUMN inv_vin TYPE CHAR(6) USING inv_vin::char(6);
ALTER TABLE bodegas ALTER COLUMN numeracion SET DEFAULT 'S/N';
ALTER TABLE bodegas ALTER COLUMN fecha_registro TYPE TIMESTAMPTZ USING fecha_registro::timestamptz;
ALTER TABLE bodegas ALTER COLUMN email_institucional DROP NOT NULL;
ALTER TABLE bodegas ADD CONSTRAINT cuit_valido CHECK (cuit ~ '^[0-9]{11}$');
ALTER TABLE bodegas ADD CONSTRAINT telefono_valido CHECK (telefono ~ '^[0-9]+$');
ALTER TABLE bodegas ADD CONSTRAINT email_valido CHECK (email_institucional like '%@%');

-- Cambios en tabla cuentas
ALTER TABLE cuenta RENAME TO cuentas;
CREATE TYPE tipo_cuenta AS ENUM ('BODEGA', 'ADMINISTRADOR_APP');
ALTER TABLE cuentas ALTER COLUMN tipo TYPE tipo_cuenta USING tipo::tipo_cuenta;
ALTER TABLE cuentas ALTER COLUMN email_login TYPE VARCHAR(150);
ALTER TABLE cuentas ALTER COLUMN fecha_registro TYPE TIMESTAMPTZ USING fecha_registro::timestamptz;
ALTER TABLE cuentas ADD CONSTRAINT cuentas_email_unique UNIQUE (email_login);
ALTER TABLE cuentas ADD CONSTRAINT cuentas_bodega_unique UNIQUE (id_bodega);
ALTER TABLE cuentas ADD CONSTRAINT cuentas_bodega_fk FOREIGN KEY (id_bodega) REFERENCES bodegas (id_bodega);
ALTER TABLE cuentas ADD CONSTRAINT cuenta_bodega_requerida CHECK (
  (tipo = 'BODEGA' AND id_bodega IS NOT NULL) OR
  (tipo = 'ADMINISTRADOR_APP' AND id_bodega IS NULL)
);

-- Cambios en tabla responsables
ALTER TABLE responsable RENAME TO responsables;
ALTER TABLE responsables DROP CONSTRAINT IF EXISTS responsables_pkey;
ALTER TABLE responsables ADD COLUMN id_cuenta INTEGER NOT NULL;
ALTER TABLE responsables DROP COLUMN id_bodega;
ALTER TABLE responsables ADD COLUMN fecha_baja TIMESTAMPTZ;
ALTER TABLE responsables ALTER COLUMN fecha_registro TYPE TIMESTAMPTZ USING fecha_registro::timestamptz;
ALTER TABLE responsables ALTER COLUMN dni TYPE VARCHAR(8);
ALTER TABLE responsables ADD CONSTRAINT responsables_pkey PRIMARY KEY (id_responsable);
ALTER TABLE responsables ADD CONSTRAINT responsables_cuenta_fk FOREIGN KEY (id_cuenta) REFERENCES cuentas (id_cuenta);
ALTER TABLE responsables ADD CONSTRAINT dni_valido CHECK (dni IS NULL OR dni ~ '^[0-9]{7,8}$');
CREATE UNIQUE INDEX un_responsable_activo_por_cuenta ON responsables (id_cuenta) WHERE activo = true;
