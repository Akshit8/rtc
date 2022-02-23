let ws = new WebSocket("ws://localhost:8080");

ws.send("hi from client");

ws.onmessage = (message) => {
  console.log(`Recieved from server: ${message.data}`);
};
