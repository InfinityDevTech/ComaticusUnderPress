function parse(raw) {
    return JSON.parse(raw)
  }

module.exports = class SocketHandler{
    constructor(client, id, headers, words) {
        this.client = client;
        this.id = id
        this.ip = headers.ip
        this.words = words
        this.hasWord = false
        this.pickWord()
        console.log(`New client connected! ID: ${this.id} IP: ${this.ip}`)
    }

    pickWord() {
        console.log(this.words.length)
        this.curWord = this.words[Math.floor(Math.random() * this.words.length)]
        this.hasWord = true
        console.log(this.curWord)
        this.client.send(JSON.stringify({type: "word", word: this.curWord}))
    }
}