const morphdom = require("morphdom");
const cave = require("../index");
// import morphdom from "morphdom";
describe("Array", () => {
  describe("#indexOf()", () => {
    it("should return -1 when the value is not present", () => {
      assert.equal(-1, [1, 2, 3].indexOf(4));
    });
  });
  describe("document stuff", () => {
    it("should be able to access the document", () => {
      const testCases = [
        [
          "<div></div>",
          "<div></div><div></div>",
          [{ t: 0, d: "\u003cdiv\u003e\u003c/div\u003e", i: 2 }],
        ],
        ["<div></div><div></div>", "<div></div>", [{ t: 1, i: 2 }]],
        [
          "<div foo=bar one=two>Hello</div>",
          '<div foo="baz" three="four">Hello</div>',
          [
            {
              t: 2,
              a: [
                ["foo", "baz"],
                ["three", "four"],
              ],
              i: 1,
            },
          ],
        ],
        ["<div>Hello</div>", "<div>World</div>", [{ t: 3, d: "World", i: 2 }]],
        [
          "<div>Hello</div><div>World</div>",
          "<div>World</div><div>Hello</div>",
          [
            { t: 3, d: "World", i: 2 },
            { t: 3, d: "Hello", i: 4 },
          ],
        ],
        [
          "<div>Hello</div>",
          '<span foo="bar">Hello</span>',
          [
            {
              t: 4,
              d: '\u003cspan foo="bar"\u003eHello\u003c/span\u003e',
              i: 1,
            },
          ],
        ],
      ];
      for (let i = 0; i < testCases.length; i++) {
        const data = testCases[i];
        document.querySelector("body").innerHTML = data[0];
        let node = document.querySelector("body").cloneNode(true);
        cave.apply(node, data[2]);
        morphdom(document.querySelector("body"), node);
        document.querySelector("body").innerHTML.should.equal(data[1]);
      }
      console.log(document.querySelector("body"));
    });
  });
});
