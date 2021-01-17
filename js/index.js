const PATCH_TYPE_INSERT = 0;
const PATCH_TYPE_REMOVE = 1;
const PATCH_TYPE_ATTRIBUTES = 2;
const PATCH_TYPE_TEXT = 3;
const PATCH_TYPE_ELEMENT = 4;
const _apply = (node, patches, index) => {
  console.log(patches, index);
  if (patches.length == 0) {
    return [0, null];
  }
  patch = patches[0];
  if (patch.i === index) {
    patches = patches.slice(1);
    if (patch.t === PATCH_TYPE_REMOVE) {
      node.parentNode.removeChild(node);
    } else if (patch.t == PATCH_TYPE_ATTRIBUTES) {
      while (node.attributes.length > 0)
        node.removeAttribute(node.attributes[0].name);
      for (let i = 0; i < patch.a.length; i++) {
        const attr = patch.a[i];
        node.setAttribute(attr[0], attr[1]);
      }
    } else if (patch.t === PATCH_TYPE_TEXT) {
      console.log(patch.d, node);
      node.data = patch.d;
    } else if (patch.t === PATCH_TYPE_ELEMENT) {
      node.replaceWith(stringToNode(patch.d));
    }
  }
  if (patches.length == 0) {
    return [0, null];
  }
  patch = patches[0];
  if (patch.t === PATCH_TYPE_INSERT && index + 1 === patch.i) {
    node.parentNode.appendChild(stringToNode(patch.d));
    patches = patches.slice(1);
  }
  for (let i = 0; i < node.childNodes.length; i++) {
    const child = node.childNodes[i];
    [index, patches] = _apply(child, patches, index + 1);
  }
  return [index, patches];
};
const apply = (node, patches) => {
  _apply(node, patches, 0);
};

const stringToNode = (input) => {
  let doc = new DOMParser().parseFromString(input, "text/html");
  return doc.body.firstChild;
};

module.exports = { apply: apply };
