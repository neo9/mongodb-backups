const path = "./seed.json";
const dbName = "integration"
const collectionName = "test"


db = db.getSiblingDB(dbName);

print(db[collectionName].countDocuments());