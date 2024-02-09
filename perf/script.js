import http from 'k6/http';
import grpc from 'k6/net/grpc';
import * as YAML from "k6/x/yaml";
import { check } from 'k6';
import { b64encode } from 'k6/encoding';
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// 1 - Upload wasm file
// 2 - Upload schema
// 3 - Upload the job from test data


export let options = {
  vus: 30,
  iterations: 100
};

const clientCtl = new grpc.Client();
const clientRepo = new grpc.Client();
clientCtl.load(['C:\\Users\\Andres\\projects\\go\\jobico\\internal\\api\\proto'], 'ctl.proto');
clientRepo.load(['C:\\Users\\Andres\\projects\\go\\jobico\\internal\\api\\proto'], 'repo.proto');
const jobyml = open('../internal/test/testdata/job.yml')
const wasmFile = open('../internal/test/testdata/echo.wasm', 'b')
const schemaFile = open('../internal/test/testdata/schema.json', 'b')

export function setup() {
  clientCtl.connect('localhost:50052', {
    plaintext: true
  });
  clientRepo.connect('localhost:50053', {
    plaintext: true
  });

  const gettenant = {
    ID: "tenant_1",
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
          name: "sch1",
          type: 1
        },
        tenant: "tenant_1"
      }
    }
    const wasm = {
      tenantFile: {
        file: {
          content: b64encode(wasmFile),
          name: "run1",
          type: 2
        },
        tenant: "tenant_1"
      }
    }
    const tenant = {
      tenant: {
        ID: "tenant_1",
        Name: "tenant_1"
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
  const response = http.post('http://localhost:8080/events/tenant_1/event_id_1', JSON.stringify(payload), {
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
