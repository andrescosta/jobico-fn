/* 
0 --T0
    Creates a new P0 package

T1                              T2
1 ---  Sleep for 5 seconds       Sends package for P0
...
5 ---  Creates a new package P1  Sends package for P0
       Sleep for 5 seconds
...
10 --- Creates a new package P2  Sends package for P0
       Sleep for 5 seconds

(
  T1: Continue uploading a new package until stops.
  T2: Continue sending events for P0 until stops.
)
*/


import { Test } from './lib/test.js'
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import * as config from './configs.js'
import { sleep } from 'k6'

const test = new Test(config.TENANT, config.HOST_CTL, config.HOST_LISTENER, config.HOST_REPO, config.TLS)
test.LoadFileBin('../internal/test/testdata/echo.wasm')
test.LoadFileBin('../internal/test/testdata/schema.json')
test.LoadFile('../internal/test/testdata/job.yml')

export const options = {
  scenarios: {
    tenant_sending_events: {
      executor: 'constant-vus',
      env: { WORKFLOW: 'sending_events' },
      vus: 10,
      duration: '60s',
    },
    tenant_adding_pkg: {
      executor: 'per-vu-iterations',
      exec: 'addPackage',
      env: { WORKFLOW: 'adding' },
      startTime: '1s',
      vus: 3,
      iterations: 5,
      maxDuration: '20s',
    },
  },
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
  test.Close();
}

export default () => {
  test.SendEventV1Random()
}

export function addPackage() {
  if (__ITER === 0) {
    test.Connect();
  }
  const id = randomIntBetween(0, 1000)
  const jobid = 'job_id_' + id
  const jobname = 'job_id_name_' + id
  const idq = 'queue_id_' + id
  const nameq = 'queue_name_' + id
  test.AddPackageFileForJobWithTemplate(jobid, jobname, idq, nameq, '../internal/test/testdata/job.yml');
  sleep(5)
}

export function teardown() {
  test.Close()
}
