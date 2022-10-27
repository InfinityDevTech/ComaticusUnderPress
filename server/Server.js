const express = require('express')
const WebSocket = require('ws')

const servHandler = require('./SocketStuffs')

const sockServ = new WebSocket.WebSocketServer({ noServer: true })
new servHandler(sockServ)

const app = express()
const server = app.listen(80)

server.on('upgrade', (request, socket, head) => {
    sockServ.handleUpgrade(request, socket, head, socket => {
      sockServ.emit('connection', socket, request, request.headers);
    });
  });