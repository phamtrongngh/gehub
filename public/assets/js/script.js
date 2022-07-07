let socket;

function connect() {
  const alias = String(document.getElementById("alias-txt").value);
  const port = Number(document.getElementById("port-txt").value);

  if (!port || port < 1024 || port > 65535) {
    showNotiBox("Port must be in range 1024-65535", "error");
    return;
  }

  if (!/^[\w-]{0,30}$/.test(alias)) {
    showNotiBox(
      "Alias must be words, numbers, hyphens or underscores (less than 30 characters)",
      "error"
    );
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

  socket.on("forward", async (data) => {
    const { path, method, port, body, headers } = data;
    let response, resStatus, resHeaders, resBody;

    try {
      response = await fetch(`http://localhost:${port}/${path}`, {
        headers,
        method,
        body,
      });
      resStatus = response.status;

      resHeaders = {};
      response.headers.forEach(function (value, key) {
        resHeaders[key] = value;
      });

      const contentType = resHeaders["content-type"];
      if (contentType && contentType.indexOf("application/json") !== -1) {
        resBody = await response.json();
      } else {
        resBody = await response.text();
      }
    } catch (e) {
      showNotiBox(`
        Cannot forward request to local server. 
        Ensure that the server has been enabled CORS and accepts all HTTP request methods.
      `, "error");
      resHeaders = {};
      resStatus = 502;
      resBody = {};
    }

    insertLog(resStatus, method, path);

    socket.emit("forward", {
      status: resStatus,
      headers: resHeaders,
      body: resBody,
    });
  });

  socket.on("disconnect", disconnect);
}

function disconnect(msg) {
  hideElement(connectResult);
  showElement(loading);
  showNotiBox("Connection closed: " + msg, "error")

  if (socket.connected) {
    socket.disconnect();
  }

  clearLogs();

  hideElement(loading);
  showElement(connectForm);
}

connectBtn.addEventListener("click", connect);
disconnectBtn.addEventListener("click", disconnect);
