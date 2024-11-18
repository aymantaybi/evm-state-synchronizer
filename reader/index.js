const { Level } = require("level");
const { Account, BN, bufferToHex, rlp, keccak256, toBuffer } = require("ethereumjs-util");
const { SecureTrie: Trie } = require("merkle-patricia-tree");

const db = new Level("../../chaindata");

db.on("open", () => {
  console.log("Database is open and ready for operations.");
});

db.on("error", (err) => {
  console.error("Error occurred while opening the database:", err);
});

const stateRoot = "0xdd2b02b747fc61ac1d9a586324287a06f940dfbb370070f0f176685ac63c9029"; // 40021727
const stateRootBuffer = Buffer.from(stateRoot.slice(2), "hex");

const trie = new Trie(db, stateRootBuffer);

(async function () {
  trie
    .createReadStream()
    .on("data", console.log)
    .on("end", () => {
      console.log("End.");
    });
})();
