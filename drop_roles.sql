DROP USER sutartys;
REVOKE ALL ON sutartis, brezinys, leidimas FROM sutartys_roles;
REVOKE ALL ON sutartis_id_seq FROM sutartys_roles;
DROP ROLE sutartys_roles;

DROP USER breziniai;
REVOKE ALL ON sutartis, brezinys, leidimas FROM breziniai_roles;
REVOKE ALL ON brezinys_id_seq FROM breziniai_roles;
DROP ROLE breziniai_roles;

DROP USER leidimai;
REVOKE ALL ON brezinys, leidimas FROM leidimai_roles;
REVOKE ALL ON leidimas_id_seq FROM leidimai_roles;
DROP ROLE leidimai_roles;
