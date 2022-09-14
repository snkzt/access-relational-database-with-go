DROP TABLE IF EXISTS album;
CREATE TABLE album (
  id        INT AUTO_INCREMENT NOT NULL,
  title     VARCHAR(128) NOT NULL,
  artist    VARCHAR(255) NOT NULL,
  price     DECIMAL(5,2) NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO album
  (title, artist, price)
VALUES
  ('Mack The Knife', 'Ella Fitzgerald', 56.99),
  ('Ain't Got No, I Got Life', 'Nina Simone', 63.99),
  ('Rhumba Azul', 'Nat King Cole', 77.99),
  ('Today', 'Johnny Hartman', 54.98);
