CREATE TABLE public_sources.kfm03 (
                                      id bigserial NOT NULL,
                                      num bigserial NOT NULL,
                                      decision varchar(200) NULL,
                                      weapon_sanction_listed_at date NULL,
                                      information varchar(4000) NULL,
                                      CONSTRAINT kfm03_pk PRIMARY KEY (id)
);

CREATE TABLE public_sources.practiceXML (
    id bigserial NOT NULL PRIMARY KEY,
    naming varchar(100) NULL,
    survived_years int NULL,
    second_naming varchar(100) NULL,
    ip varchar(100) NULL,
    dog varchar(100) null,
    naznayu varchar(100)

)