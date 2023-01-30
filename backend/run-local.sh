cd cmd/server
echo "Building server binary"
go build
cd ../..
mv cmd/server/server server-binary

export MIGRATIONS_PATH='./resources/db/sqlite'
export DB_FILENAME='db.sql'
export ENABLE_CREATE_USERS=true
export JWT_SECRET='d40e532c75861f10567379d723d96f74'

./server-binary
