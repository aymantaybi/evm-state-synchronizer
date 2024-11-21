const { Level } = require("level");
const { Account, BN, bufferToHex, rlp, keccak256, toBuffer } = require("ethereumjs-util");
const { SecureTrie: Trie } = require("merkle-patricia-tree");

const db = new Level("../../roninchain/ronin/chaindata");

db.on("open", () => {
  console.log("Database is open and ready for operations.");
});

db.on("error", (err) => {
  console.error("Error occurred while opening the database:", err);
});

const stateRoot = "0xddc8f1b241f9090547501d92c2a943a41e8b076f14f2836be2cd8b4b1f6053c4"; // 40079306
const stateRootBuffer = Buffer.from(stateRoot.slice(2), "hex");

const trie = new Trie(db, stateRootBuffer);

const address = "0x0b7007c13325c48911f73a2dad5fa5dcbf808adc"; // usdc
const addressBuffer = Buffer.from(address.slice(2), "hex");

(async function () {
  const data = await trie.get(addressBuffer);
  const acc = Account.fromRlpSerializedAccount(data);

  console.log("-------State-------");
  console.log(`nonce: ${acc.nonce}`);
  console.log(`balance in wei: ${acc.balance}`);
  console.log(`storageRoot: ${bufferToHex(acc.stateRoot)}`);
  console.log(`codeHash: ${bufferToHex(acc.codeHash)}`);

  const storageTrie = trie.copy();
  storageTrie.root = acc.stateRoot;

  console.log("------Storage------");
  const stream = storageTrie.createReadStream();
  stream
    .on("data", (data) => {
      console.log(`key: ${bufferToHex(data.key)}`);
      console.log(`Value: ${bufferToHex(rlp.decode(data.value))}`);
    })
    .on("end", () => {
      console.log("Finished reading storage.");
    });
})();
