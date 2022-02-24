let ws = new WebSocket("ws://localhost:8080/chat");

ws.send("hi from client");

ws.onmessage = (message) => {
  console.log(`Recieved from server: ${message.data}`);
};

ws.close(1011, "Going away");