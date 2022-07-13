docker build . -t microblogging-db
docker run -p 54321:5432 microblogging-db
#postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable