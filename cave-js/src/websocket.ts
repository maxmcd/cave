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
    this.caveSubmitListener = this.caveSubmitListener.bind(this);
    this.caveClickListener = this.caveClickListener.bind(this);
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
    this.send(new ClientMessage(MessageType.Init, componentIDs).serialize());
  }
  send(data: any): void {
    console.log("sending msg to server", data);
    return this.websocket.send(data);
  }
  onmessage(e: MessageEvent): any {
    console.log(e.data);
    let msg = new ServerMessage(JSON.parse(e.data));
    console.log(msg);
    if (msg.event === MessageType.Patch) {
      let patches: Array<Patch> = msg.data.map((e) => diff.expandPatch(e));
      diff.apply(document.querySelector("[cave-component]"), patches);
      this.addEventListeners();
    }
    if (msg.event === MessageType.Error) {
      throw new Error("Server error: " + msg.data[0]);
    }
  }

  addEventListeners() {
    // TODO: only add listeners to new code?
    // how expensive is this?
    console.log("adding event listeners");
    document.querySelectorAll("[cave-submit]").forEach((elem) => {
      elem.addEventListener("submit", this.caveSubmitListener);
    });
    document.querySelectorAll("[cave-click]").forEach((elem) => {
      elem.addEventListener("click", this.caveClickListener);
    });
  }
  findElementContext(
    element: HTMLElement
  ): [string | undefined, string | undefined] {
    let componentID =
      element.closest("[cave-component]")?.getAttribute("cave-component") ||
      undefined;
    let subcomponentID =
      element
        .closest("[cave-subcomponent]")
        ?.getAttribute("cave-subcomponent") || undefined;
    return [componentID, subcomponentID];
  }
  caveClickListener(e: Event): void {
    let element = <HTMLElement>e.target;
    const [componentID, subcomponentID] = this.findElementContext(element);
    if (!componentID) {
      // if we're not in a cave component we shouldn't act
      return;
    }
    e.preventDefault();
    let name =
      element.getAttribute("cave-click") ||
      element.closest("[cave-click]")?.getAttribute("cave-click") ||
      undefined;
    this.send(
      new ClientMessage(
        MessageType.Click,
        [],
        componentID,
        name,
        subcomponentID
      ).serialize()
    );
  }
  caveSubmitListener(e: Event): void {
    console.log("got submit", e, this);
    let formElement = <HTMLFormElement>e.target;

    let name = formElement.getAttribute("cave-submit");
    const [componentID, subcomponentID] = this.findElementContext(formElement);
    if (!componentID) {
      // if we're not in a cave component we shouldn't act
      return;
    }
    e.preventDefault();
    let formMessage = new SubmitMessage(
      componentID,
      name || "",
      formElement,
      subcomponentID
    );
    this.send(formMessage.serialize());
  }
}

const wsScheme = (): string => {
  if (window.location.protocol == "http:") {
    return "ws";
  }
  return "wss";
};
