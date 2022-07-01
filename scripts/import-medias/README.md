# Import Medias

1. Download all datasets
```bash
wget https://datasets.imdbws.com/name.basics.tsv.gz
wget https://datasets.imdbws.com/title.akas.tsv.gz
wget https://datasets.imdbws.com/title.basics.tsv.gz
wget https://datasets.imdbws.com/title.crew.tsv.gz
wget https://datasets.imdbws.com/title.episode.tsv.gz
wget https://datasets.imdbws.com/title.principals.tsv.gz
wget https://datasets.imdbws.com/title.ratings.tsv.gz
```
2. Unzip all archives ``gzip -d *.tsv.gz``
3. Move all **tsv** files to the medias folder
4. Run the container
```
docker build -t hypertube-import-medias .
docker run hypertube-import-medias
```
