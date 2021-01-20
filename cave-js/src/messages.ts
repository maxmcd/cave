export type WebsocketMessage = [string,  Array<any>, string];
export enum MessageType {
  Patch = "p",
  Init = "init",
  Error = "error",
  Submit = "submit",
}

export class ClientMessage {
  type: MessageType
  data: Array<any>;
  componentID?: string;
  name?: string;
  subcomponentID?: string;
  constructor(
    type: MessageType,
    data: Array<any>,
    componentID?: string,
    name?: string,
    subcomponentID?: string,
  ) {
    this.type = type;
    this.data = data;
    this.componentID  = componentID;
    this.name = name;
    this.subcomponentID = subcomponentID;
  }
  prepare(): Array<any> {
    let out: Array<any> = [this.type, this.data]
    if (!this.componentID) {
      return out
    }
    out.push(this.componentID)
    if (!this.name) {
      return out
    }
    out.push(this.name)
    if (!this.subcomponentID) {
      return out
    }
    out.push(this.subcomponentID)
    return out
  }
  serialize(): string {
    return JSON.stringify(this.prepare());
  }
}

export class ServerMessage {
  componentID: string;
  event: string;
  data: Array<any>;
  constructor(input: WebsocketMessage) {
    this.event = input[0];
    this.data = input[1];
    this.componentID = input[2];
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
    super(MessageType.Submit, [formDataMap], componentID, name);
  }
}
