// Load the seed data from the seed.json file
console.log("Deleting database records");
console.log("-------------------------------------------");

const path = "./seed.json";
const dbName = "integration"
const collectionName = "test"


db = db.getSiblingDB(dbName);

//Start the DDBB anew
db.getCollectionNames().forEach((collectionName) => {
    db[collectionName].deleteMany({});
});

console.log("Finished");
