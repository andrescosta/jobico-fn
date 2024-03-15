import { Test } from './lib/test.js'
import * as config from './configs.js'

const test = new Test(config.TENANT, config.HOST_CTL, config.HOST_LISTENER, config.HOST_REPO, config.TLS)
test.LoadFileBin('../internal/test/testdata/echo.wasm')
test.LoadFileBin('../internal/test/testdata/schema.json')
test.LoadFile('../internal/test/testdata/job.yml')

export let options = {
  vus: 1,
  iterations: 1
};

export function setup() {
  test.Connect();
  const e = test.ExistsTenant();
  if (!e) {
    test.AddTenant();
    test.UploadWasmFile('run1', '../internal/test/testdata/echo.wasm');
    test.UploadSchemaFile('sch1', '../internal/test/testdata/schema.json');
    test.AddPackageWithFile('../internal/test/testdata/job.yml');
  }
  test.Close()
}

export default () => {
  test.SendEventV1Random()
}

export function teardown() {
}
