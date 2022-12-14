cd cmd/server
echo "Building server binary"
go build
cd ../..
mv cmd/server/server server-binary

export MIGRATIONS_PATH='./resources/db/sqlite'
export DB_FILENAME='db.sql'

./server-binary