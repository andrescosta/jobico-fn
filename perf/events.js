import { Test } from './lib/test.js'

const HOST_CTL = 'ctl:443'
const HOST_REPO = 'repo:443'
const HOST_LISTENER = 'listener'
const TENANT = 'tenant_1'
const test = new Test(TENANT, HOST_CTL, HOST_LISTENER, HOST_REPO)
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
