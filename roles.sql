CREATE USER sutartys WITH PASSWORD 'sutartys';
REVOKE ALL ON brezinys, leidimas, sutartis, brezinys_id_seq, leidimas_id_seq, sutartis_id_seq FROM sutartys;
CREATE ROLE sutartys_roles WITH LOGIN PASSWORD 'sutartys';
REVOKE ALL ON brezinys, leidimas, sutartis, brezinys_id_seq, leidimas_id_seq, sutartis_id_seq FROM sutartys_roles;
GRANT UPDATE, INSERT ON sutartis TO sutartys_roles;
GRANT SELECT, DELETE ON sutartis, brezinys, leidimas TO sutartys_roles;
GRANT ALL ON sutartis_id_seq TO sutartys_roles;
GRANT sutartys_roles TO sutartys;

CREATE USER breziniai WITH PASSWORD 'breziniai';
REVOKE ALL ON brezinys, leidimas, sutartis, brezinys_id_seq, leidimas_id_seq, sutartis_id_seq FROM breziniai;
CREATE ROLE breziniai_roles WITH LOGIN PASSWORD 'breziniai';
REVOKE ALL ON brezinys, leidimas, sutartis, brezinys_id_seq, leidimas_id_seq, sutartis_id_seq FROM breziniai_roles;
GRANT UPDATE, INSERT ON brezinys TO breziniai_roles;
GRANT SELECT ON sutartis TO breziniai_roles;
GRANT SELECT, DELETE ON brezinys, leidimas TO breziniai_roles;
GRANT ALL ON brezinys_id_seq TO breziniai_roles;
GRANT breziniai_roles TO breziniai;

CREATE USER leidimai WITH PASSWORD 'leidimai';
REVOKE ALL ON brezinys, leidimas, sutartis, brezinys_id_seq, leidimas_id_seq, sutartis_id_seq FROM leidimai;
CREATE ROLE leidimai_roles WITH LOGIN PASSWORD 'leidimai';
REVOKE ALL ON brezinys, leidimas, sutartis, brezinys_id_seq, leidimas_id_seq, sutartis_id_seq FROM leidimai_roles;
GRANT UPDATE, INSERT ON leidimas TO leidimai_roles;
GRANT SELECT ON brezinys TO leidimai_roles;
GRANT SELECT, DELETE ON leidimas TO leidimai_roles;
GRANT ALL ON leidimas_id_seq TO leidimai_roles;
GRANT leidimai_roles TO leidimai;
