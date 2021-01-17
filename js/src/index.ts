enum PatchType {
  Insert = 0,
  Remove,
  Attributes,
  Text,
  Element,
}

interface Patch {
  data: string;
  index: number;
  attributes: Array<Array<string>>;
  type: number;
}

interface PatchOnTheWire {
  d: string;
  i: number;
  a: Array<Array<string>>;
  t: number;
}

const expandPatch = (p: PatchOnTheWire): Patch => {
  return {
    data: p.d,
    index: p.i,
    attributes: p.a,
    type: p.t,
  };
};

const apply = (node: HTMLElement, patches: Array<Patch>): void => {
  _apply(node, patches, 0);
};

const _apply = (
  node: any,
  patches: Array<Patch>,
  index: number
): [Array<Patch>, number] => {
  if (patches.length == 0) {
    return [null, 0];
  }
  let patch = patches[0];

  if (patch.index === index) {
    patches = patches.slice(1);
    if (patch.type === PatchType.Remove) {
      node.parentNode.removeChild(node);
    } else if (patch.type == PatchType.Attributes) {
      while (node.attributes.length > 0)
        node.removeAttribute(node.attributes[0].name);
      for (let i = 0; i < patch.attributes.length; i++) {
        const attr = patch.attributes[i];
        node.setAttribute(attr[0], attr[1]);
      }
    } else if (patch.type === PatchType.Text) {
      node.data = patch.data;
    } else if (patch.type === PatchType.Element) {
      node.replaceWith(stringToNode(patch.data));
    }
  }

  if (patches.length == 0) {
    return [null, 0];
  }
  patch = patches[0];
  if (patch.type === PatchType.Insert && index + 1 === patch.index) {
    node.parentNode.appendChild(stringToNode(patch.data));
    patches = patches.slice(1);
  }
  for (let i = 0; i < node.childNodes.length; i++) {
    const child = node.childNodes[i];
    [patches, index] = _apply(child, patches, index + 1);
  }
  return [patches, index];
};

const stringToNode = (input: string): ChildNode => {
  let doc = new DOMParser().parseFromString(input, "text/html");
  return doc.body.firstChild;
};

module.exports = { apply, expandPatch };
