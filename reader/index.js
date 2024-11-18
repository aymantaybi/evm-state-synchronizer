const { Level } = require("level");
const { Account, BN, bufferToHex, rlp, keccak256, toBuffer } = require("ethereumjs-util");
const { BaseTrie: Trie } = require("merkle-patricia-tree");

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

const address = "0xc1eb47de5d549d45a871e32d9d082e7ac5d2e3ed";
const addressBuffer = keccak256(toBuffer(address));

(async function () {
  trie
    .createReadStream()
    .on("data", console.log)
    .on("end", () => {
      console.log("End.");
    });
})();
