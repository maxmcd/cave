export type WebsocketMessage = [string, string, Array<any>];
export enum MessageType {
  Patch = "p",
  Init = "init",
  Error = "error",
}

export class ClientMessage {
  componentID: string;
  event: string;
  name: string;
  data: Array<any>;
  constructor(
    componentID: string,
    event: string,
    name: string,
    data: Array<any>
  ) {
    this.componentID = componentID;
    this.event = event;
    this.name = name;
    this.data = data;
  }
  serialize(): string {
    return JSON.stringify([this.componentID, this.event, this.name, this.data]);
  }
}

export class ServerMessage {
  componentID: string;
  event: string;
  data: Array<any>;
  constructor(input: WebsocketMessage) {
    this.componentID = input[0];
    this.event = input[1];
    this.data = input[2];
  }
  serialize(): string {
    return JSON.stringify([this.componentID, this.event, this.data]);
  }
}

export class SubmitMessage extends ClientMessage {
  data: [string, Record<string, string>];
  constructor(componentID: string, name: string, form: HTMLFormElement) {
    let formData = new FormData(form);
    let formDataMap: Record<string, string> = {};
    formData.forEach((v, k) => {
      // TODO: file support
      formDataMap[k] = v.toString();
    });
    super(componentID, "submit", name, [formDataMap]);
  }
}
