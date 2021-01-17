module.exports = function (config) {
  config.set({
    frameworks: ["mocha", "chai", "browserify"],
    files: ["test/**/*.js"],
    reporters: ["progress"],
    port: 9876, // karma web server port
    colors: true,
    logLevel: config.LOG_INFO,
    browsers: ["ChromeHeadless"],
    preprocessors: {
      "test/**/*.js": ["browserify"],
    },
    autoWatch: false,
    // singleRun: false, // Karma captures browsers, runs the tests and exits
    concurrency: Infinity,
  });
};
