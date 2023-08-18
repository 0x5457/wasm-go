// Usage: node --experimental-wasi-unstable-preview1 setup_suite.js

const path = require('path')
const fs = require('fs/promises')
const { WASI } = require('wasi')

const testNames = ['address', 'block', 'i32', 'i64', 'f32', 'f64'];
async function main() {
  const wast2jsonMod = await WebAssembly.compileStreaming(fetch('https://registry-cdn.wapm.io/contents/wasmer/wabt/1.0.37/out/wast2json.wasm'))
  for (const testName of testNames) {
    console.log(`downloading testcase ${testName}...`)
    const testUrl = `https://raw.githubusercontent.com/WebAssembly/testsuite/main/${testName}.wast`
    const wast = await fetch(testUrl).then(res => res.text())
    const wastPath = path.join(__dirname, `suite/`)
    await makeDir(wastPath)
    await fs.writeFile(path.join(wastPath, `${testName}.wast`), wast)
    const wasi = new WASI({
      version: 'preview1',
      args: ['wast2json', `/suite/${testName}.wast`, '-o', `/suite/json/${testName}.json`],
      preopens: {
        '/suite': path.join(__dirname, 'suite'),
      },
    });
    const inst = await WebAssembly.instantiate(wast2jsonMod, wasi.getImportObject())

    await makeDir(path.join(__dirname, 'suite/json'))
    console.log(`running wast2json...`)
    wasi.start(inst);
  }
}

async function makeDir(p) {
  try { await fs.mkdir(p, { recursive: true }) } catch (e) { }
}

main()
