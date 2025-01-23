// Load the seed data from the seed.json file
console.log("Starting the initialisation of the database");

console.log("-------------------------------------------")
const fs = require('fs');
const path = "/tmp/data/seed.json";
const dbName = "integration"
const collectionName = "test"

const seedData = JSON.parse(fs.readFileSync(path, 'utf8'));

db = db.getSiblingDB(dbName);

//Start the DDBB anew
db.getCollectionNames().forEach((collectionName) => {
    db[collectionName].drop();
});

db.createCollection(collectionName);
db[collectionName].insertMany(seedData);

console.log("Finished");
