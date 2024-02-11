import { Test } from './test.js'

const HOST_CTL = 'localhost:50052'
const HOST_REPO = 'localhost:50053'
const HOST_LISTENER = 'localhost:8080'
const TENANT = 'tenant_1'
const test = new Test(TENANT, HOST_CTL, HOST_LISTENER, HOST_REPO)

export let options = {
  vus: 1,
  iterations: 1
};

export function setup() {
  test.Init();

  const e = test.ExistsTenant();
  if (e) {
    test.AddTenant(TENANT, TENANT);
    test.UploadWasmFile('run1', '../internal/test/testdata/echo.wasm');
    test.UploadSchemaFile('sch1', '../internal/test/testdata/schema.json');
    test.AddPackageFile('../internal/test/testdata/job.yml');
  }
}

export default () => {
  test.SendEventV1Random(TENANT, randomFirstName, randomLastName, age)
}

export function teardown() {
  test.Close()
}
