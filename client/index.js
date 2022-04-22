function generateUserID() {
  let letters = "abcdefghigklmnopqrstuvyxwzABCDEFGHIGKLMNOPQRSTUVYXWZ";
  let digits = "1234567890";
  let symbols = letters + digits;
  let userID = "";
  for (let i = 0; i < 16; i++) {
    userID += symbols[Math.floor(Math.random() * (symbols.length - 1))];
  }
  return userID;
}

let user_id;
if (localStorage.getItem("user_id") === null) {
  user_id = generateUserID();
  localStorage.setItem("user_id", user_id);
} else {
  user_id = localStorage.getItem("user_id");
}

connect();

function connect() {
  let socket = new WebSocket("ws://localhost:4221/ws");

  const params = new Proxy(new URLSearchParams(window.location.search), {
    get: (searchParams, prop) => searchParams.get(prop),
  });

  function changeStatus(status) {
    let statusElement = document.getElementById("status");
    let color = status ? "green" : "red";
    statusElement.innerHTML = `<span style="color: ${color}"> ${
      status ? "online" : "offline"
    }<span>`;
  }

  changeStatus(false);

  socket.onopen = function (e) {
    console.log("[open] Соединение установлено");

    socket.send(
      JSON.stringify({
        type: "emit",
        event: "user connected",
        data: {
          local_user_id: user_id,
          room_id: params.room_id,
        },
      })
    );

    changeStatus(true);
  };

  socket.onmessage = function (event) {
    let message = JSON.parse(event.data);
    console.log(message);
    switch (message.type) {
      case "emit":
        let clients = document.getElementById("clients");
        switch (message.event) {
          case "join user":
            let client = document.createElement("li");
            client.innerText =
              message.data.id + " : " + message.data.local_user_id;
            clients.appendChild(client);
            break;
          case "your id":
            let socketId = document.getElementById("socketId");
            socketId.innerHTML = message.data.id;
            break;
          case "leave user":
            for (var i = 0; i < clients.childNodes.length; i++) {
              if (
                clients.childNodes[i].innerText.split(" : ")[0] ===
                message.data.id
              ) {
                clients.removeChild(clients.childNodes[i]);
              }
            }
            break;
        }
        break;
    }
  };

  socket.onclose = function (event) {
    changeStatus(false);
    if (event.wasClean) {
      console.log(
        `[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`
      );
    } else {
      console.log("[close] Соединение прервано");
    }
    setTimeout(function () {
      console.log("[reopen] Переподключение");
      connect();
    }, 5000);
  };

  socket.onerror = function (error) {
    console.log(`[error] ${error.message}`);
  };
}
