const fs = require("fs");
const path = require("path");
const rm = require("rimraf");
const mkdirp = require("mkdirp");
const sh = require("shelljs");

const VERSION_CMD = sh.exec("git describe --abbrev=0 --tags");
const VERSION_DATE_CMD = sh.exec("go run ./scripts/date.go");

const VERSION = VERSION_CMD.replace("\n", "").replace("\r", "");
const VERSION_DATE = VERSION_DATE_CMD.replace("\n", "").replace("\r", "");

const ROOT = __dirname;
const TEMPLATES = path.join(ROOT, "templates");

async function updateTranExtension(ghTranDir) {
  const templatePath = path.join(TEMPLATES, "gh-tran");
  const template = fs.readFileSync(templatePath).toString("utf-8");

  const templateReplaced = template
    .replace("CLI_VERSION", VERSION)
    .replace("CLI_VERSION_DATE", VERSION_DATE);

  fs.writeFileSync(path.join(ghTranDir, "gh-tran"), templateReplaced);
}

async function updateExtension() {
  const tmp = path.join(__dirname, "tmp");
  const extensionDir = path.join(tmp, "gh-tran");

  mkdirp.sync(tmp);
  rm.sync(extensionDir);

  console.log(`cloning https://github.com/abdfnx/gh-tran to ${extensionDir}`);

  sh.exec(`git clone https://github.com/abdfnx/gh-tran.git ${extensionDir}`)

  console.log(`done cloning abdfnx/gh-tran to ${extensionDir}`);

  console.log("updating local git...");

  await updateTranExtension(extensionDir);
}

updateExtension().catch((err) => {
  console.error(`error running scripts/gh-tran/gh-trn.js`, err);
  process.exit(1);
});
