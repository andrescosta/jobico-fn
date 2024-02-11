import grpc from 'k6/net/grpc';
import * as YAML from "k6/x/yaml";
import http from 'k6/http';
import { b64encode } from 'k6/encoding';

export function Api(hostCtl, hostListener, hostRepo) {
  this.hostCtl = hostCtl
  this.hostListener = hostListener
  this.hostRepo = hostRepo
  this.clientCtl = new grpc.Client();
  this.clientRepo = new grpc.Client();
  this.clientCtl.load(['../internal/api/proto'], 'ctl.proto');
  this.clientRepo.load(['../internal/api/proto'], 'repo.proto');
}

Api.prototype.Init = function () {
  this.clientCtl.connect(this.hostCtl, {
    plaintext: true
  });
  this.clientRepo.connect(this.hostRepo, {
    plaintext: true
  });
}

Api.prototype.GetTenant = function (id) {
  const tenant = {
      ID: id,
  }
  console.log(tenant)
  const response = this.clientCtl.invoke('/Control/Tenants', tenant);
  return response
}

Api.prototype.AddTenant = function (id, name) {
  const tenant = {
    tenant: {
      ID: id,
      Name: name
    }
  }
  var response = this.clientCtl.invoke('/Control/AddTenant', tenant);
  return response
}

Api.prototype.UploadWasmFile = function (tenant, fileId, path) {
  const wasm = {
    tenantFile: {
      file: {
        content: this.File(path),
        name: fileId,
        type: 2
      },
      tenant: tenant
    }
  }

  const response = this.clientRepo.invoke('Repo/AddFile', wasm);
  return response
}

Api.prototype.UploadSchemaFile = function (tenant, fileId, path) {
  const schema = {
    tenantFile: {
      file: {
        content: this.File(path),
        name: fileId,
        type: 1
      },
      tenant: tenant
    }
  }

  const response = this.clientRepo.invoke('Repo/AddFile', schema);
  return response

}

Api.prototype.AddPackageFile = function (path) {
  const jobyml = open(path)
  const job = YAML.parse(jobyml)
  response = clientCtl.invoke('Control/AddPackage', job);

}

Api.prototype.SendEvent = function (tenant, evt, evtBody) {
  const url = 'http://' + this.hostListener + '/events/' + tenant + '/' + evt;

  const response = http.post(url, JSON.stringify(evtBody), {
    headers: { 'Content-Type': 'application/json' },
  });
  return response
}

Api.prototype.File = function (path) {
  const file = open(path, 'b');
  return b64encode(file);
}

Api.prototype.Close = function(){
  this.clientRepo.close();
  this.clientCtl.close();
}