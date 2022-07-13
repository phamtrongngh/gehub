const socket = io("", { transports: ["websocket"] });

socket.on("forward", async (data) => {
  const { path, method, port, body, headers } = data;
  let response, resStatus, resHeaders, resBody;

  try {
    response = await fetch(`http://0.0.0.0:${port}/${path}`, {
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

    // if content-type is image, then encode it to base64 and insert into resBody
    if (resHeaders["content-type"].find((item) => item.includes("image"))) {
      const imageBlob = await response.blob();
      const reader = new FileReader();
      reader.readAsDataURL(imageBlob);
      reader.onloadend = () => {
        const base64data = reader.result;
        resBody = base64data;
        resHeaders["content-length"] = [`${base64data.length}`];

        socket.emit("forward", {
          status: resStatus,
          headers: resHeaders,
          body: resBody,
        });
      };
    } else {
      resBody = await response.text();

      socket.emit("forward", {
        status: resStatus,
        headers: resHeaders,
        body: resBody,
      });
    }
  } catch (e) {
    showNotiBox(
      `
      Cannot forward request to local server. Ensure that the server 
      has been enabled CORS and accepts all HTTP request methods.`,
      "error"
    );
    resHeaders = {};
    resStatus = 502;
    resBody = "";

    socket.emit("forward", {
      status: resStatus,
      headers: resHeaders,
      body: resBody,
    });
  }
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
  showNotiBox("Expose successfully", "success");
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
