window.onload = function () {};

function sendID(e) {
  if (e.key === "Enter") {
    s.setName(e.target.value);
    console.log(s.getName());
    e.target.value = "";
    enableDisable("#user_id", true);
    enableDisable("#message", false);
    connect(s.getName());
  }
}

function enableDisable(name, state) {
  document.querySelector(name).disabled = state;
}

function connect(name) {
  let url = "ws://localhost:8888";
  if (name) {
    url += "?name=" + name;
  }

  const ws = new WebSocket(url);
  ws.onopen = (e) => {
    // console.log("ws opened");
    ws.send("onopen message from client");
  };

  ws.onmessage = (e) => {
    // console.log("ws receive message");
    console.log(e.data);
  };

  ws.onclose = (e) => {
    console.log("ws close");
    console.log(e);
  };

  ws.onerror = (err) => {
    console.log("ws error");
    console.log(err);
  };
}

function store() {
  let _name = "";
  let nameState = true;
  let messageState = false;

  function getName() {
    return _name;
  }
  function setName(name) {
    _name = name;
  }

  return {
    setName,
    getName,
  };
}
window.s = store();
