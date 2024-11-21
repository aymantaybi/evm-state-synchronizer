const { Level } = require("level");
const { Account, BN, bufferToHex, rlp, keccak256, toBuffer } = require("ethereumjs-util");
const { SecureTrie: Trie } = require("merkle-patricia-tree");

const db = new Level("../roninchain/ronin/chaindata");

db.on("open", () => {
  console.log("Database is open and ready for operations.");
});

db.on("error", (err) => {
  console.error("Error occurred while opening the database:", err);
});

const stateRoot = "0xddc8f1b241f9090547501d92c2a943a41e8b076f14f2836be2cd8b4b1f6053c4"; // 40079306
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
