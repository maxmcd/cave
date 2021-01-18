const morphdom = require("morphdom");
import type { Patch } from "./diff";
import diff from "./diff";
import {
  SubmitMessage,
  ServerMessage,
  MessageType,
  ClientMessage,
} from "./messages";

export class WebsocketHandler {
  websocket: WebSocket;
  constructor() {
    this.addEventListeners();
    this.websocket = new WebSocket(
      `${wsScheme()}://${window.location.host}${
        window.location.pathname
      }?cavews`
    );
    this.websocket.onmessage = this.onmessage.bind(this);
    this.websocket.onopen = this.onopen.bind(this);
  }
  onopen(e: Event) {
    let componentIDs: Array<string> = [];
    let components = document.querySelectorAll("[cave-component]");
    for (let i = 0; i < components.length; i++) {
      const component = components.item(i);
      const val = component.getAttribute("cave-component");
      if (val) {
        componentIDs.push(val);
      }
    }
    this.websocket.send(
      new ClientMessage("", MessageType.Init, "", componentIDs).serialize()
    );
  }
  onmessage(e: MessageEvent): any {
    console.log(e.data);
    let msg = new ServerMessage(JSON.parse(e.data));
    if (msg.event === MessageType.Patch) {
      let patches: Array<Patch> = msg.data.map((e) => diff.expandPatch(e));
      // is this what we want? is cloning expensive here?
      let node = document.querySelector("[cave-component]")?.cloneNode(true);
      diff.apply(node, patches);
      morphdom(document.querySelector("[cave-component]"), node);
    }
    if (msg.event === MessageType.Error) {
      throw new Error("Server error: " + msg.data[0]);
    }
  }

  addEventListeners() {
    console.log("adding event listeners");
    document.querySelectorAll("[cave-submit]").forEach((elem) => {
      let name = elem.getAttribute("cave-submit");
      if (!name) {
        // attributes must have a value
        return;
      }
      elem.addEventListener("submit", this.caveSubmitListener(<string>name));
    });
  }

  caveSubmitListener(name: string): (e: Event) => void {
    return (e) => {
      let formElement = <HTMLFormElement>e.currentTarget;
      let componentID = formElement
        .closest("[cave-component]")
        ?.getAttribute("cave-component");
      if (!componentID) {
        // if we're not in a cave component we shouldn't act
        return;
      }
      e.preventDefault();
      let formMessage = new SubmitMessage(componentID, name, formElement);
      this.websocket.send(formMessage.serialize());
    };
  }
}

const wsScheme = (): string => {
  if (window.location.protocol == "http:") {
    return "ws";
  }
  return "wss";
};
