build for docker alpine GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

restore database:
pg_restore -j 8 -U user -d db /path/srbolabdb.dump.20230718.220901.pgdmp