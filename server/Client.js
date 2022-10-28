function parse(raw) {
    return JSON.parse(raw)
  }

module.exports = class SocketHandler{
    constructor(client, id, headers, words, server) {
        this.client = client;
        this.id = id
        this.ip = headers.ip
        this.server = server
        this.words = words
        this.hasWord = false
        this.pickWord()
        setInterval(() => {client.send(JSON.stringify({type: 'heartbeat'}))}, 2000)
        console.log(`New client connected | ID: ${this.id} | IP: ${this.ip} | Word: ${this.curWord}`)
        this.checkGuess()
    }

    checkGuess() {
        this.client.on('message', (data) => {
            let parsed = parse(data)
            if (parsed.Type == "guess") {
                if (parsed.Guess == this.curWord) {
                    this.client.send(JSON.stringify({type: "correct"}))
                } else {
                    this.client.send(JSON.stringify({type: "incorrect"}))
                    this.server.clients.forEach((client) => {
                            client.send(JSON.stringify({type: "leaked", ip: this.ip}));
                    })
                }
            }
        })
    }

    pickWord() {
        this.curWord = this.words[Math.floor(Math.random() * this.words.length)]
        this.hasWord = true
        this.client.send(JSON.stringify({type: "word", word: this.curWord}))
    }
}