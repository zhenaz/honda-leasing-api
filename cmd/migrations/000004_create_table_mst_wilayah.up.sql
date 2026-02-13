-- Schema: mst 

DROP TABLE IF EXISTS mst.province CASCADE;
DROP TABLE IF EXISTS mst.kabupaten CASCADE;
DROP TABLE IF EXISTS mst.kecamatan CASCADE;
DROP TABLE IF EXISTS mst.kelurahan CASCADE;
DROP TABLE IF EXISTS mst.locations CASCADE;

-- Fokus Provinsi: Jawa Barat saja
-- 1. province <<mst>>
CREATE TABLE mst.province (
    prov_id    BIGSERIAL PRIMARY KEY,
    prov_name  VARCHAR(85) UNIQUE NOT NULL
);

-- 2. kabupaten <<mst>>
CREATE TABLE mst.kabupaten (
    kab_id     BIGSERIAL PRIMARY KEY,
    kab_name   VARCHAR(85) NOT NULL,
    prov_id    BIGINT NOT NULL REFERENCES mst.province(prov_id) ON DELETE RESTRICT
);

-- 3. kecamatan <<mst>>
CREATE TABLE mst.kecamatan (
    kec_id     BIGSERIAL PRIMARY KEY,
    kec_name   VARCHAR(85) NOT NULL,
    kab_id     BIGINT NOT NULL REFERENCES mst.kabupaten(kab_id) ON DELETE RESTRICT
);

-- 4. kelurahan <<mst>>
CREATE TABLE mst.kelurahan (
    kel_id     BIGSERIAL PRIMARY KEY,
    kel_name   VARCHAR(85) NOT NULL,
    kec_id     BIGINT NOT NULL REFERENCES mst.kecamatan(kec_id) ON DELETE RESTRICT
);

-- 5. locations <<mst>>
CREATE TABLE mst.locations (
    location_id    BIGSERIAL PRIMARY KEY,
    street_address TEXT,
    postal_code    VARCHAR(10),
    longitude      numeric(9,6),
    latitude       numeric(9,6),
    kel_id         BIGINT REFERENCES mst.kelurahan(kel_id) ON DELETE SET NULL
);

-- =============================================
-- INDEX untuk performa query 
-- =============================================
CREATE INDEX idx_kabupaten_prov ON mst.kabupaten(prov_id);
CREATE INDEX idx_kecamatan_kab  ON mst.kecamatan(kab_id);
CREATE INDEX idx_kelurahan_kec  ON mst.kelurahan(kec_id);
CREATE INDEX idx_locations_kel  ON mst.locations(kel_id);

-- =============================================
-- Seed Data (INSERT raw data)
-- =============================================

-- 1. province: hanya Jawa Barat
INSERT INTO mst.province (prov_name) VALUES ('Jawa Barat');

-- 2. kabupaten: semua kab/kota di Jawa Barat 
INSERT INTO mst.kabupaten (kab_name, prov_id) VALUES
('Kabupaten Bandung', 1),
('Kabupaten Bandung Barat', 1),
('Kabupaten Bekasi', 1),
('Kabupaten Bogor', 1),
('Kabupaten Ciamis', 1),
('Kabupaten Cianjur', 1),
('Kabupaten Cirebon', 1),
('Kabupaten Garut', 1),
('Kabupaten Indramayu', 1),
('Kabupaten Karawang', 1),
('Kabupaten Kuningan', 1),
('Kabupaten Majalengka', 1),
('Kabupaten Pangandaran', 1),
('Kabupaten Purwakarta', 1),
('Kabupaten Subang', 1),
('Kabupaten Sukabumi', 1),
('Kabupaten Sumedang', 1),
('Kabupaten Tasikmalaya', 1),
('Kota Bandung', 1),
('Kota Bekasi', 1),
('Kota Bogor', 1),
('Kota Cimahi', 1),
('Kota Cirebon', 1),
('Kota Depok', 1),
('Kota Sukabumi', 1),
('Kota Tasikmalaya', 1),
('Kota Banjar', 1);

-- 3. kecamatan: contoh 5 kecamatan di Kota Bandung (kab_id asumsi 19 = Kota Bandung)
INSERT INTO mst.kecamatan (kec_name, kab_id) VALUES
('Andir', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Antapani', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Arcamanik', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Astanaanyar', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Bandung Kulon', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung'));

-- 4. kelurahan: contoh kelurahan di Kecamatan Bandung Kulon 
INSERT INTO mst.kelurahan (kel_name, kec_id) VALUES
('Caringin', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cibuntu', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cijerah', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Gempolsari', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cigondewah Kaler', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cigondewah Kidul', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Warung Muncang', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cigondewah Rahayu', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon'));

-- 5. locations: alamat dummy di Bandung
INSERT INTO mst.locations (street_address, postal_code, longitude, latitude, kel_id) VALUES
('Jl. Soekarno Hatta No. 123, Caringin', '40212', '107.612345', '-6.912345', (SELECT kel_id FROM mst.kelurahan WHERE kel_name = 'Caringin')),
('Jl. Terusan Buahbatu No. 45, Cijerah', '40213', '107.598765', '-6.923456', (SELECT kel_id FROM mst.kelurahan WHERE kel_name = 'Cijerah')),
('Komplek Griya Bandung Indah, Warung Muncang', '40211', '107.585678', '-6.935678', (SELECT kel_id FROM mst.kelurahan WHERE kel_name = 'Warung Muncang'));