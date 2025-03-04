
db.createUser(
    {
      user: "godbuser",
      pwd: "godbpass",
      roles: [ { role: "readWrite", db: "godb" } ]
    }
);

db.createCollection('assets');
db.createCollection('audienceCharacteristics');
db.createCollection('userFavorites');
db.createCollection('users');
