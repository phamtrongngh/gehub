const connectForm = document.querySelector("#connect-form");
const connectBtn = connectForm.querySelector("#connect-btn");

const connectResult = document.querySelector("#connect-result");
const disconnectBtn = connectResult.querySelector("#disconnect-btn");
const requestTable = connectResult.querySelector("#request-table");

const loading = document.querySelector("#loading");

const notiBox = document.querySelector("#noti-box");

const copyables = document.querySelectorAll(".copyable");
copyables.forEach((c) => {
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

function createBadgeElement(text, theme) {
  let element;
  element = document.createElement("span");
  element.className = "badge";
  switch (theme) {
    case "red":
      element.className += " badge-red";
      break;
    case "blue":
      element.className += " badge-blue";
      break;
    case "green":
      element.className += " badge-green";
      break;
    case "orange":
      element.className += " badge-orange";
      break;
    default:
      break;
  }
  element.textContent = text;
  return element;
}

function insertLog(status, method, path) {
  const newRow = document.createElement("tr");
  newRow.innerHTML = "<td></td><td></td><td></td><td></td>";

  var date = new Date(),
    fmtDate =
      [date.getDate(), date.getMonth() + 1, date.getFullYear()].join("/") +
      " " +
      [date.getHours(), date.getMinutes(), date.getSeconds()].join(":");
  newRow.cells[0].appendChild(createBadgeElement(fmtDate));

  switch (true) {
    case status >= 200 && status <= 299:
      newRow.cells[1].appendChild(createBadgeElement(status, "green"));
      break;
    case status >= 300 && status <= 399:
      newRow.cells[1].appendChild(createBadgeElement(status, "blue"));
      break;
    case status >= 400 && status <= 499:
      newRow.cells[1].appendChild(createBadgeElement(status, "orange"));
      break;
    case status >= 500 && status <= 599:
      newRow.cells[1].appendChild(createBadgeElement(status, "red"));
      break;
    default:
      newRow.cells[1].appendChild(createBadgeElement(status));
      break;
  }

  switch (method) {
    case "GET":
      newRow.cells[2].appendChild(createBadgeElement(method, "green"));
      break;
    case "POST":
      newRow.cells[2].appendChild(createBadgeElement(method, "blue"));
      break;
    case "PUT":
    case "PATCH":
      newRow.cells[2].appendChild(createBadgeElement(method, "orange"));
      break;
    case "DELETE":
      newRow.cells[2].appendChild(createBadgeElement(method, "red"));
      break;

    default:
      newRow.cells[2].appendChild(createBadgeElement(method));
      break;
  }

  newRow.cells[3].appendChild(createBadgeElement(`/${path}`));

  requestTable.tBodies[0].appendChild(newRow);
}

function clearLogs() {
  requestTable.tBodies[0].innerHTML = "";
}

function showNotiBox(msg, theme) {
  notiBox.textContent = msg;
  notiBox.className = theme + " show";

  setTimeout(function () {
    notiBox.className = notiBox.className.replace("show", "");
  }, 5000);
}
