import grpc from 'k6/net/grpc';
import * as YAML from "k6/x/yaml";
import http from 'k6/http';
import { b64encode } from 'k6/encoding';

export function Api(hostCtl, hostListener, hostRepo, tls) {
  this.tls = tls
  this.hostCtl = hostCtl
  this.hostListener = hostListener
  this.hostRepo = hostRepo
  this.clientCtl = new grpc.Client();
  this.clientRepo = new grpc.Client();
  this.clientCtl.load(['../internal/api/proto'], 'ctl.proto');
  this.clientRepo.load(['../internal/api/proto'], 'repo.proto');
  this.files = new Map();
}


Api.prototype.Connect = function () {
  this.clientCtl.connect(this.hostCtl, {
    plaintext: !this.tls
  });
  this.clientRepo.connect(this.hostRepo, {
    plaintext: !this.tls
  });
}

Api.prototype.LoadFile = function(path){
    this.files.set(path,loadFileUTF(path))
}

Api.prototype.LoadFileBin = function(path){
  this.files.set(path,b64encode(loadFileBin(path)))
}

Api.prototype.File = function(path){
  return this.files.get(path)
}


Api.prototype.GetTenant = function (id) {
  const tenant = {
      ID: id,
  }
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
Api.prototype.GetPackageObj = function(path){
  const jobyml = this.File(path)
  const job = YAML.parse(jobyml)
  return job
}

Api.prototype.AddPackageFile = function (path) {
  const job = this.GetPackageObj(path)  
  const response = this.InvokeAddPackage(job);
  return response
}

Api.prototype.InvokeAddPackage = function (pkg) {
  const response = this.clientCtl.invoke('Control/AddPackage', pkg);
  return response
}


Api.prototype.SendEvent = function (tenant, evt, evtBody) {
  const url = this.hostListener + '/events/' + tenant + '/' + evt;

  const response = http.post(url, JSON.stringify(evtBody), {
    headers: { 'Content-Type': 'application/json' },
  });
  return response
}

Api.prototype.Close = function(){
  this.clientRepo.close();
  this.clientCtl.close();
}


function loadFileBin (path) {
  const file = open(path, 'b');
  return file;
}

function loadFileUTF (path) {
  const file = open(path);
  return file;
}
