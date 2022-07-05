let socket;

const connectForm = document.querySelector("#connect-form");
const connectBtn = connectForm.querySelector("#connect-btn");
const connectResult = document.querySelector("#connect-result");
const disconnectBtn = connectResult.querySelector("#disconnect-btn");
const loading = document.querySelector("#loading");

document.querySelectorAll(".copyable").forEach((c) => {
  c.addEventListener("click", (e) => {
    navigator.clipboard.writeText(e.target.textContent);
  });
});

function hideElement(element) {
  element.style.display = "none";
}

function showElement(element) {
  element.style.display = "block";
}

function connect() {
  const alias = String(document.getElementById("alias-txt").value);
  const port = Number(document.getElementById("port-txt").value);

  if (!port || port < 1024 || port > 65535) {
    console.log("port invalid");
    return;
  }

  if (!/^[\w-]{0,30}$/.test(alias)) {
    console.log("alias invalid");
    return;
  }

  hideElement(connectForm);
  showElement(loading);

  socket = io("", {
    transports: ["websocket"],
    reconnection: false,
    query: {
      port,
      alias,
    },
  });

  socket.on("connect", () => {
    socket.emit("info");
  });

  socket.on("info", (data) => {
    hideElement(loading);

    connectResult.querySelector(
      "#forward-url"
    ).textContent = `${window.location.origin}/${data.alias}`;

    connectResult.querySelector(
      "#local-url"
    ).textContent = `http://localhost:${data.port}`;
    
    showElement(connectResult);

  });

  socket.on("forward", (data) => {
    console.log("forward", data);
  });

  socket.on("disconnect", (data) => {
    console.log("disconnect", data);
  });
}

function disconnect() {
  hideElement(connectResult);
  showElement(loading);

  socket.disconnect();

  hideElement(loading);
  showElement(connectForm);
}

connectBtn.addEventListener("click", connect);
disconnectBtn.addEventListener("click", disconnect);
