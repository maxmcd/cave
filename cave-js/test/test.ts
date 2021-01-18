const morphdom = require("morphdom");
import diff from "../src/diff";
import type { PatchOnTheWire } from "../src/";

describe("Applying diffs", () => {
  it("should be able to apply the diffs", () => {
    // prettier-ignore
    const testCases: Array<[string, string, Array<PatchOnTheWire>]> = [
      [
        "<div></div>",
        "<div></div><div></div>",
        [{ t: 0, d: "\u003cdiv\u003e\u003c/div\u003e", i: 2 }],
      ],
      [
        "<div></div><div></div>",
        "<div></div>",
        [{ t: 1, i: 2 }]],
      [
        "<div foo=bar one=two>Hello</div>",
        '<div foo="baz" three="four">Hello</div>',
        [{t: 2, a: [["foo", "baz"], ["three", "four"]], i: 1}],
      ],
      [
        "<div>Hello</div>",
        "<div>World</div>",
        [{ t: 3, d: "World", i: 2 }]],
      [
        "<div>Hello</div><div>World</div>",
        "<div>World</div><div>Hello</div>",
        [{ t: 3, d: "World", i: 2 }, { t: 3, d: "Hello", i: 4 }],
      ],
      [
        "<div>Hello</div>",
        '<span foo="bar">Hello</span>',
        [{t: 4, d: '\u003cspan foo="bar"\u003eHello\u003c/span\u003e', i: 1}],
      ],
    ];
    for (let i = 0; i < testCases.length; i++) {
      const data = testCases[i];
      getBody().innerHTML = data[0];
      let node = getBody().cloneNode(true);
      node;

      diff.apply(
        node,
        data[2].map((p) => diff.expandPatch(p))
      );
      morphdom(getBody(), node);
      getBody().innerHTML.should.equal(data[1]);
    }
  });
});

function getBody(): HTMLBodyElement {
  const out = document.querySelector("body");
  if (!out) {
    throw Error("unable to query for body");
  }
  return out;
}
