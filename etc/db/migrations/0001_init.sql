

CREATE TABLE IF NOT EXISTS materials (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	label TEXT NOT NULL,
	class TEXT NOT NULL,
	description TEXT NOT NULL,
	special BOOLEAN NOT NULL DEFAULT false,
	UNIQUE (label)
	);
	COMMENT ON TABLE materials IS 'Table storing information about filament materials';
	COMMENT ON COLUMN materials.label IS 'The detailed name of the material (e.g. Basic PLA, Hyper PLA, PETG, TPU 85A PLA+ etc.)';
	COMMENT ON COLUMN materials.class IS 'The material class/type (e.g. PLA, PETG, TPU, ABS, etc.)';
	COMMENT ON COLUMN materials.description IS 'The common name of the material (e.g. PLA, Hyper PLA, PETG, TPU 85A PLA+, etc.)';
	COMMENT ON COLUMN materials.special IS 'Is this a special material (e.g. composite filaments with wood, metal, or carbon fiber additives)?';

CREATE TABLE IF NOT EXISTS locations (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	label TEXT NOT NULL,
	description TEXT NOT NULL,
	capacity INTEGER NOT NULL DEFAULT 1,
	printable BOOLEAN NOT NULL DEFAULT false,
	tally BOOLEAN NOT NULL DEFAULT true,
	UNIQUE (label)
	);
	COMMENT ON TABLE locations IS 'Table storing information about filament storage locations';
	COMMENT ON COLUMN locations.label IS 'The name of the storage location (e.g. Shelf A, Box 1, Drawer 3, etc.)';
	COMMENT ON COLUMN locations.description IS 'A more detailed explanation of the storage location (e.g. "Top shelf in the closet", "Plastic box under the desk", etc.)';
	COMMENT ON COLUMN locations.capacity IS 'The maximum number of spools that can be stored in this location';
	COMMENT ON COLUMN locations.printable IS 'Whether this storage location is printable (i.e. is attached to or part of a 3D printer)';
	COMMENT ON COLUMN locations.tally IS 'Should this location be included in spool tallies?';

CREATE TABLE IF NOT EXISTS stores (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	label TEXT NOT NULL,
	url TEXT,
	UNIQUE (label)
	);
	COMMENT ON TABLE stores IS 'Table storing information about stores where filament was purchased';
	COMMENT ON COLUMN stores.label IS 'The name of the store (e.g. Prusament, Amazon, etc.)';
	COMMENT ON COLUMN stores.url IS 'The URL of the store''s website (e.g. https://www.ebay.com/)';

CREATE TABLE IF NOT EXISTS brands (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	label TEXT NOT NULL,
	active BOOLEAN NOT NULL DEFAULT true,
	store_id BIGINT NULL REFERENCES stores(id),
	UNIQUE (label)
	);
	COMMENT ON TABLE brands IS 'Table storing information about filament brands';
	COMMENT ON COLUMN brands.label IS 'The brand name (e.g. Prusament, Anycubic, Sunlu, etc.)';
	COMMENT ON COLUMN brands.active IS 'Whether or not the brand is still in business';
	COMMENT ON COLUMN brands.store_id IS 'Does the brand have a primary store where their filament is sold?';

CREATE TABLE IF NOT EXISTS colors (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	label TEXT NOT NULL,
	hex_code TEXT NOT NULL,
	alias TEXT,
	UNIQUE (label, hex_code)
	);
	COMMENT ON TABLE colors IS 'Table storing information about filament colors';
	COMMENT ON COLUMN colors.label IS 'The name of the color (e.g. "Red", "Blue", "Transparent", etc.)';
	COMMENT ON COLUMN colors.hex_code IS 'The hexadecimal color code (e.g. "#FF0000" for red)';
	COMMENT ON COLUMN colors.alias IS 'An alternative name for the color (e.g. "Crimson Red" for "#DC143C")';

CREATE TABLE IF NOT EXISTS spools (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	material_id BIGINT NOT NULL REFERENCES materials(id),
	brand_id BIGINT NOT NULL REFERENCES brands(id),
	location_id BIGINT NOT NULL REFERENCES locations(id),
	store_id BIGINT NOT NULL REFERENCES stores(id),
	weight NUMERIC(10,2) NOT NULL DEFAULT 0,
	combined_weight NUMERIC(10,2) NOT NULL DEFAULT 0,
	current_weight NUMERIC(10,2) NOT NULL DEFAULT 0,
	price NUMERIC(10,4) NOT NULL DEFAULT 0,
	empty BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
	);
	COMMENT ON TABLE spools IS 'Table storing information about individual filament spools';
	COMMENT ON COLUMN spools.weight IS 'The marketed weight of the just filament on the spool in grams';
	COMMENT ON COLUMN spools.combined_weight IS 'The initial, packaged weight of the both the filament and the spool in grams';
	COMMENT ON COLUMN spools.current_weight IS 'The current weight of the both the filament and the spool in grams';
	COMMENT ON COLUMN spools.price IS 'The cost of the filament at the time of purchase';
	COMMENT ON COLUMN spools.empty IS 'Whether the spool is empty'; -- different from deleted_at (soft delete), which indicates the record shouldn't have existed
	COMMENT ON COLUMN spools.created_at IS 'The timestamp when the spool record was created';
	COMMENT ON COLUMN spools.updated_at IS 'The timestamp when the spool record was last updated';
	COMMENT ON COLUMN spools.deleted_at IS 'The timestamp when the spool record was deleted (soft delete)';

CREATE TABLE IF NOT EXISTS spool_colors (
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	spool_id BIGINT NOT NULL REFERENCES spools(id),
	color_id BIGINT NOT NULL REFERENCES colors(id)
	);
	COMMENT ON TABLE spool_colors IS 'Associative table linking spools to their colors (many-to-many relationship)';

---- create above / drop below ----

DROP TABLE IF EXISTS spool_colors CASCADE;
DROP TABLE IF EXISTS spools CASCADE;
DROP TABLE IF EXISTS colors CASCADE;
DROP TABLE IF EXISTS brands CASCADE;
DROP TABLE IF EXISTS stores CASCADE;
DROP TABLE IF EXISTS locations CASCADE;
DROP TABLE IF EXISTS materials CASCADE;
