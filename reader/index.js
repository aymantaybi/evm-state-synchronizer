const { Level } = require("level");
const { Account, BN, bufferToHex, rlp } = require("ethereumjs-util");
const { SecureTrie: Trie } = require("merkle-patricia-tree");

const db = new Level("./chaindata");

const stateRoot = "0xdd2b02b747fc61ac1d9a586324287a06f940dfbb370070f0f176685ac63c9029"; // 40021727
const stateRootBuffer = Buffer.from(stateRoot.slice(2), "hex");

const trie = new Trie(db, stateRootBuffer);

const address = "0xc1eb47de5d549d45a871e32d9d082e7ac5d2e3ed";
const addressBuffer = Buffer.from(address.slice(2), "hex");

(async function () {
  const data = await trie.get(addressBuffer);
  const acc = Account.fromAccountData(data);

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
