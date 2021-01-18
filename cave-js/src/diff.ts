const morphdom = require("morphdom");

enum PatchType {
  Insert = 0,
  Remove,
  Attributes,
  Text,
  Element,
}

export interface Patch {
  data: string;
  index: number;
  attributes: Array<Array<string>>;
  type: number;
}

export interface PatchOnTheWire {
  d?: string;
  i: number;
  a?: Array<Array<string>>;
  t: number;
}

const expandPatch = (p: PatchOnTheWire): Patch => {
  let patch: Patch = {
    index: p.i,
    type: p.t,
    data: "",
    attributes: [],
  };

  if (p.d) {
    patch.data = p.d;
  }
  if (p.a) {
    patch.attributes = p.a;
  }
  return patch;
};

function printNode(node: any, level: number) {
  console.log("\t".repeat(level), node.nodeName);
  node.childNodes.forEach((node) => {
    printNode(node, level + 1);
  });
}

const apply = (node: any, patches: Array<Patch>): void => {
  // printNode(node, 0);
  _apply(node, patches, 0);
};

const _apply = (
  node: HTMLElement,
  patches: Array<Patch>,
  index: number
): [Array<Patch>, number] => {
  if (patches.length == 0) {
    return [[], 0];
  }
  let patch = patches[0];

  if (patch.index === index) {
    patches = patches.slice(1);
    if (patch.type === PatchType.Remove) {
      console.log("removing node", patch, node);
      let parent = node.parentNode;
      if (parent) {
        parent.removeChild(node);
      } // error otherwise?
    } else if (patch.type == PatchType.Attributes) {
      while (node.attributes.length > 0)
        node.removeAttribute(node.attributes[0].name);
      for (let i = 0; i < patch.attributes.length; i++) {
        const attr = patch.attributes[i];
        node.setAttribute(attr[0], attr[1]);
      }
      console.log("updating attributes", patch, node);
    } else if (patch.type === PatchType.Text) {
      console.log("updating text", patch, node);
      if (node.nodeType == Node.TEXT_NODE) {
        let unknown = <unknown>node;
        let textNode = <Text>unknown;
        textNode.data = patch.data;
      } else {
        console.log("innerhtml", node.innerHTML, patch.data);
        node.innerHTML = patch.data;
      }
    } else if (patch.type === PatchType.Element) {
      console.log("replacing element", patch, node);
      morphdom(node, patch.data);
    }
  }

  if (patches.length == 0) {
    return [[], 0];
  }
  patch = patches[0];
  if (patch.type === PatchType.Insert && index + 1 === patch.index) {
    console.log("inserting element", patch, node);
    let parent = node.parentNode;
    if (parent) {
      let patchNode = stringToNode(patch.data);
      console.log("patchNode", patchNode);
      if (patchNode) {
        parent.append(patchNode);
      }
    }
    patches = patches.slice(1);
  }
  for (let i = 0; i < node.childNodes.length; i++) {
    const child = node.childNodes[i];
    [patches, index] = _apply(<HTMLElement>child, patches, index + 1);
  }
  return [patches, index];
};

const stringToNode = (input: string): ChildNode | null => {
  let doc = new DOMParser().parseFromString(input, "text/html");
  if (!doc.body.firstChild) {
    return document.createTextNode(input);
  }
  return doc.body.firstChild;
};

export default { apply, expandPatch };
