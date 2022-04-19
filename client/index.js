let socket = new WebSocket("ws://localhost:8080/ws");

function changeStatus(status) {
    let statusElement = document.getElementById("status");
    let color = status ? 'green' : 'red';
    statusElement.innerHTML = `<span style="color: ${color}"> ${status ? 'online' : 'offline'}<span>`;
}

changeStatus(false);

socket.onopen = function(e) {
    console.log("[open] Соединение установлено");
    changeStatus(true);
};

socket.onmessage = function(event) {
    let message = JSON.parse(event.data);
    if (message.type == 'emit') {
        console.log('event ', message.event);
        if (message.event === 'new user') {
            let clients = document.getElementById('clients');

            let client = document.createElement("li");
            client.innerText = message.data.id;

            clients.appendChild(client);
        }
        if (message.event === 'your id') {
            let socketId = document.getElementById('socketId');
            socketId.innerHTML = message.data.id;
        }
        if (message.event === 'leave user') {
            let clients = document.getElementById('clients');
            for (var i = 0; i < clients.childNodes.length; i++) {
                if (clients.childNodes[i].innerText === message.data.id) {
                    clients.removeChild(clients.childNodes[i]);
                }
            }
        }
    }
};
  
socket.onclose = function(event) {
    changeStatus(false);
    if (event.wasClean) {
        console.log(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
    } else {
        console.log('[close] Соединение прервано');
    }
};

socket.onerror = function(error) {
    console.log(`[error] ${error.message}`);
};