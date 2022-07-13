const socket = io("", { transports: ["websocket"] });

function getLocalHostname(){
  let userAgent = navigator.userAgent;  
  if(userAgent.match(/chrome|chromium|crios/i)){
      return "0.0.0.0"
    }else if(userAgent.match(/firefox|fxios/i)){
      return "0.0.0.0"
    }  else if(userAgent.match(/safari/i)){
      return "127.0.0.1"
    } else if(userAgent.match(/opr\//i)){
      return "0.0.0.0";
    } else if(userAgent.match(/edg/i)){
      return "0.0.0.0"
    }else{
      return "0.0.0.0"
    }
}

socket.on("forward", async (data) => {
  const { path, method, port, body, headers } = data;
  let response, resStatus, resHeaders, resBody;

  try {
    response = await fetch(`http://${getLocalHostname()}:${port}/${path}`, {
      headers,
      method,
      body,
    });
    resStatus = response.status;

    insertLog(resStatus, method, path);

    resHeaders = {};
    response.headers.forEach(function (value, key) {
      resHeaders[key] = value.replace(/ /g, "").split(",");
    });

    resBody = await response.text();
  } catch (e) {
    showNotiBox(
      `
      Cannot forward request to local server. Ensure that the server 
      has been enabled CORS (Access-Control-Allow-Origin, Access-Control-Allow-Headers, Access-Control-Allow-Methods).`,
      "error"
    );
    resHeaders = {};
    resStatus = 502;
    resBody = "";
  }

  socket.emit("forward", {
    status: resStatus,
    headers: resHeaders,
    body: resBody,
  });
});

// Expose port
exposeBtn.addEventListener("click", () => {
  const alias = String(document.getElementById("alias-txt").value);
  const port = Number(document.getElementById("port-txt").value);

  if (!port || port < 1024 || port > 65535) {
    showNotiBox("Port must be in range 1024-65535", "error");
    return;
  }

  if (!/^[\w-]{0,30}$/.test(alias)) {
    showNotiBox(
      "Alias must be letters, numbers, hyphens or underscores (less than 30 characters)",
      "error"
    );
    return;
  }

  hideElement(exposeForm);
  showElement(loading);

  socket.emit("expose", { alias, port });
});

socket.on("expose", (data) => {
  hideElement(loading);

  if (data.error) {
    showNotiBox(data.error, "error");
    showElement(exposeForm);
    return;
  }

  forwardUrl.textContent = data.proxyPublicUrl;
  localUrl.textContent = `http://localhost:${data.port}`;
  showElement(exposeResult);
  showNotiBox(`Expose your local server successfully. Make sure that the local server has been 
    enabled CORS (Access-Control-Allow-Origin, Access-Control-Allow-Headers, 
    Access-Control-Allow-Methods)`, "success"
  );
});

// Unexpose port
unexposeBtn.addEventListener("click", () => {
  hideElement(exposeResult);
  showElement(loading);
  socket.emit("unexpose");
});

socket.on("unexpose", () => {
  clearLogs();
  hideElement(loading);
  showElement(exposeForm);
  showNotiBox("Unexpose successfully", "success");
});

// disconnect
socket.on("disconnect", (msg) => {
  clearLogs();
  hideElement(exposeResult);
  showElement(exposeForm);
  showNotiBox("Server disconnected: " + msg, "error");
});
