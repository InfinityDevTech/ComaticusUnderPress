const Client = require("./Client");
const uuid = require("uuid");
const fs = require("fs");
const readline = require('readline');


function parse(raw) {
  return JSON.parse(raw);
}

async function processLineByLine() {
  const fileStream = fs.createReadStream("words.txt");
  const words = []

  const rl = readline.createInterface({
    input: fileStream,
    crlfDelay: Infinity,
  });

  for await (const line of rl) {
    words.push(line)
  }
    return words
}

module.exports = class SocketHandler {
  constructor(server) {
    this.clients = new Map();
    this.socket = server;
    this.socket.on("connection", this.initClient.bind(this));
    this.initWords()
  }

  async initWords() {
    this.words = await processLineByLine()
  }

  initClient(socket, req, headers) {
    let id = uuid.v4();
    let client = new Client(socket, id, headers, this.words, this.socket);
    this.clients.set(id, client);
    socket.on("close", () => {
      this.clients.delete(id);
      id = null;
      client = null;
    });
  }
};
