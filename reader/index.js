const levelup = require("levelup");
const leveldown = require("leveldown");
const { BaseTrie: Trie } = require("merkle-patricia-tree");
const Account = require("ethereumjs-account").default;
const { bufferToHex, rlp, keccak256, toBuffer } = require("ethereumjs-util");
const BN = require("bn.js");
const path = require("path");

// Correctly resolve the database path relative to the current file
const dbPath = path.resolve(__dirname, "../../chaindata");

console.log({ dbPath });

// Initialize the LevelDB database using levelup and leveldown
const db = levelup(leveldown(dbPath));

const stateRoot = "0x997847c28515099e9a040dc84c560cb83bd58a1708144b9572697c884a0b58fd"; // Block 40021727
const stateRootBuffer = Buffer.from(stateRoot.slice(2), "hex");

// Initialize the trie with the database and the state root
const trie = new Trie(db, stateRootBuffer);

(async function () {
  try {
    console.log("Starting to read trie...");

    const stream = trie.createReadStream();

    stream.on("data", (data) => {
      console.log("Key:", bufferToHex(data.key));
      console.log("Value:", bufferToHex(data.value));
    });

    stream.on("end", () => {
      console.log("End.");
    });

    stream.on("error", (err) => {
      console.error("Error during read stream:", err);
    });
  } catch (err) {
    console.error("An error occurred:", err);
  }
})();
