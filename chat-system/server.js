const { Server } = require('socket.io');

const io = new Server(3000, {});

io.on("connection", (socket)=> {
  console.log("User conected! Id: " + socket.id);
  socket.on("disconnect", ()=> {
    console.log("User has left the server :(");
  })
});

console.log("Server is listnening on port 3000");
