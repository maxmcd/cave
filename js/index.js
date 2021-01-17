const diff = require("virtual-dom/diff");
const patch = require("virtual-dom/patch");
const vdomVirtualize = require("vdom-virtualize");

var VNode = require("virtual-dom/vnode/vnode");
var VText = require("virtual-dom/vnode/vtext");

var convertHTML = require("html-to-vdom")({
  VNode: VNode,
  VText: VText,
});
var html = "<div>Hello</div>";
var html2 = "<div>World</div><div>Hello</div>";
var vtree = convertHTML(html);
var vtree2 = convertHTML(html2);

let thing = diff(vtree2, vtree);
console.log(thing);
