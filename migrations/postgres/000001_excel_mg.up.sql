CREATE TABLE public_sources.kgd14 (
                                      id bigserial NOT NULL,
                                      bin varchar(12) NULL,
                                      rnn varchar(12) NULL,
                                      "name" varchar(1000) NULL,
                                      fio_payer varchar(200) NULL,
                                      fio_director varchar(200) NULL,
                                      iin_director varchar(12) NULL,
                                      rnn_director varchar(12) NULL,
                                      number_examination varchar(500) NULL,
                                      "date" date NULL,
                                      CONSTRAINT kgd14_pk PRIMARY KEY (id)
);
CREATE INDEX kgd14_bin_index ON public_sources.kgd14 USING btree (bin);
