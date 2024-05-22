CREATE TABLE sutartis (
	id SMALLSERIAL PRIMARY KEY,
	sutartis VARCHAR NOT NULL,
	kaina VARCHAR NOT NULL
);

CREATE TABLE brezinys (
	id SERIAL PRIMARY KEY NOT NULL,
	brezinys VARCHAR NOT NULL,
	fk_sutartis_id SMALLINT NOT NULL,
	FOREIGN KEY (fk_sutartis_id) REFERENCES sutartis (id)
);

CREATE TABLE leidimas (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	leidimas VARCHAR NOT NULL,
	fk_brezinys_id INT NOT NULL,
	FOREIGN KEY (fk_brezinys_id) REFERENCES brezinys (id)
);

CREATE TABLE sutartis_id_seq (
	last_value SMALLINT,
	log_cnt INT,
	is_called bool
);

CREATE TABLE brezinys_id_seq (
	last_value INT,
	log_cnt INT,
	is_called bool
);

CREATE TABLE leidimas_id_seq (
	last_value BIGINT,
	log_cnt INT,
	is_called bool
);
