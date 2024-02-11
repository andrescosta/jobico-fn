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


import http from 'k6/http';
import grpc from 'k6/net/grpc';
import * as YAML from "k6/x/yaml";
import { check } from 'k6';
import { b64encode } from 'k6/encoding';
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const HOST_CTL = 'localhost:50052'
const HOST_REPO = 'localhost:50053'
const HOST_LISTENER = 'localhost:8080'
const TENANT = 'tenant_1'
const SCHEMA = 'sch1'
const WASM = 'run1' 
const EVENT = 'event_id_1'
const URL_LISTENER = 'http://'+HOST_LISTENER+'/events/'+TENANT+'/'+EVENT

export let options = {
  vus: 30,
  iterations: 100
};

const clientCtl = new grpc.Client();
const clientRepo = new grpc.Client();
clientCtl.load(['../internal/api/proto'], 'ctl.proto');
clientRepo.load(['../internal/api/proto'], 'repo.proto');
const jobyml = open('../internal/test/testdata/job.yml')
const wasmFile = open('../internal/test/testdata/echo.wasm', 'b')
const schemaFile = open('../internal/test/testdata/schema.json', 'b')

export function setup() {
  clientCtl.connect(HOST_CTL, {
    plaintext: true
  });
  clientRepo.connect(HOST_REPO, {
    plaintext: true
  });

  const gettenant = {
    ID: TENANT,
  }
  var response = clientCtl.invoke('/Control/Tenants', gettenant);
  check(response, {
    'status is OK': (r) => r && r.status === grpc.StatusOK,
  });
  if (response.message.Tenants.length == 0) {
    const schema = {
      tenantFile: {
        file: {
          content: b64encode(schemaFile),
          name: SCHEMA,
          type: 1
        },
        tenant: TENANT
      }
    }
    const wasm = {
      tenantFile: {
        file: {
          content: b64encode(wasmFile),
          name: WASM,
          type: 2
        },
        tenant: TENANT
      }
    }
    const tenant = {
      tenant: {
        ID: TENANT,
        Name: TENANT
      }
    }
    var response = clientCtl.invoke('/Control/AddTenant', tenant);
    check(response, {
      'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    var response = clientRepo.invoke('Repo/AddFile', schema);
    check(response, {
      'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    var response = clientRepo.invoke('Repo/AddFile', wasm);
    check(response, {
      'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    const job = YAML.parse(jobyml)
    response = clientCtl.invoke('Control/AddPackage', job);
    check(response, {
      'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
  }
}

export default () => {
  const randomFirstName = randomString(8);
  const randomLastName = randomString(10);
  const age = randomIntBetween(0,99)
  const payload = {
    data:
      [
        {
          firstName: randomFirstName,
          lastName: randomLastName,
          age: age
        }
      ]
  };
  const response = http.post(URL_LISTENER, JSON.stringify(payload), {
    headers: { 'Content-Type': 'application/json' },
  })
  check(response, {
    'status is 200': (r) => r.status === 200
  });
}

export function teardown() {
  clientRepo.close();
  clientCtl.close();
}
