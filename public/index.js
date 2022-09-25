// Rooms
const rooms = [
  "general",
  "ideas",
  "projects",
  "music"
]

let username = localStorage.getItem("chatty-username")
while (!username || !username.trim()) {
  username = prompt("Enter a username")
  if (username.trim()) {
    localStorage.setItem("chatty-username", username.trim())
  }
}

const url = new URL(window.location.href);
const wsURL = (url.protocol === "https:" ? "wss://" : "ws://") + (url.host) + "/ws" + (url.pathname)
const ws = new WebSocket(wsURL)
const room = url.pathname.replace("/", "")

const form = document.querySelector("form")
const roomList = document.querySelector("ul")
const header = document.querySelector("#header")
const messageInput = document.querySelector("input")
const ty = document.querySelector(".typing")
const online = document.querySelector("#online")
const menu = document.querySelector("#rooms-menu")

menu.addEventListener("click", function () {
  document.querySelector("#rooms").classList.toggle("open")
})

let currentTyping

rooms.forEach(r => {
  const item = document.createElement("li")
  const a = document.createElement("a")
  a.setAttribute("href", "/" + r)
  a.textContent = "#" + r
  item.appendChild(a)
  roomList.appendChild(item)
})

ws.addEventListener("open", function () {
  const m = {
    type: "join",
    payload: {
      by: username,
      room: room
    }
  }
  ws.send(JSON.stringify(m))
})

ws.addEventListener("message", function (message) {
  console.log(JSON.parse(message.data))
  const data = JSON.parse(message.data)
  switch (data.type) {
    case "join":
      console.log("Join")
      join(data.payload)
      break
    case "left":
      left(data.payload)
      break
    case "message":
      add(data.payload)
      break
    case "typing":
      typing(data.payload)
      break
    case "stoptyping":
      stoptyping(data.payload)
      break
    default:
      break
  }
})

form.onsubmit = function (e) {
  e.preventDefault();
  const m = {
    type: "message",
    payload: {
      body: messageInput.value,
      by: username,
      room
    }
  }
  console.log(m)
  ws.send(JSON.stringify(m))
  messageInput.value = ""
}

function join(message) {
  const messages = document.querySelector("#messages")
  const ele = document.createElement("div")
  ele.classList.add("message", "add-message")
  ele.textContent = "🤖: " + message.by + " has joined the room"
  messages.appendChild(ele)
  online.textContent = parseInt(message.body, 10)
}

function add(message) {
  console.log("Add message: ", message)
  const messages = document.querySelector("#messages")
  const ele = document.createElement("div")
  ele.className = "message"
  ele.textContent = message.by + ": " + message.body
  messages.appendChild(ele)
}

function left(message) {
  const messages = document.querySelector("#messages")
  const ele = document.createElement("div")
  ele.classList.add("message", "add-message")
  ele.textContent = "🤖: " + message.by + " has left the room"
  messages.appendChild(ele)
  online.textContent = parseInt(message.body, 10)
}

function typing(message) {
  console.log("typing: ", message)
  if (message.by === username) return;
  ty.textContent = message.by + " is typing ..."
  currentTyping = message.by
  console.log("Currently typing: ", currentTyping)
}

function stoptyping(message) {
  console.log("stoptyping: ", message)
  if (message.by === username) return
  if (message.by !== currentTyping) return
  ty.textContent = ""
  currentTyping = null
}

messageInput.addEventListener("focusin", function () {
  console.log("Focus in")
  const m = {
    type: "typing",
    payload: {
      room,
    }
  }
  ws.send(JSON.stringify(m))
})

messageInput.addEventListener("focusout", function () {
  console.log("Focus out")
  const m = {
    type: "stoptyping",
    payload: {
      room,
    }
  }
  ws.send(JSON.stringify(m))
})