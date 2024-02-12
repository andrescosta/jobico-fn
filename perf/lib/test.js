import { check } from 'k6';
import { Api } from './api.js'
import grpc from 'k6/net/grpc';
import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export function Test(tenant, hostCtl, hostListener, hostRepo) {
    this.api = new Api(hostCtl, hostListener, hostRepo)
    this.tenant = tenant
}

Test.prototype.LoadFile = function (path) {
    this.api.LoadFile(path)
}

Test.prototype.LoadFileBin = function (path) {
    this.api.LoadFileBin(path)
}

Test.prototype.Connect = function () {
    this.api.Connect()
}
Test.prototype.SendEventV1Random = function () {
    this.SendEventV1(randomString(8), randomString(10), randomIntBetween(0, 99))
}

Test.prototype.SendEventV1 = function (firstName, lastName, age) {
    const payload = {
        data:
            [
                {
                    firstName: firstName,
                    lastName: lastName,
                    age: age
                }
            ]
    };
    const response = this.api.SendEvent(this.tenant, 'event_id_1', payload);
    check(response, {
        'status is 200': (r) => r.status === 200
    });
}
Test.prototype.ExistsTenant = function () {
    const response = this.api.GetTenant(this.tenant);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
    return response.message.Tenants.length > 0
}

Test.prototype.GetTenant = function () {
    const response = this.api.GetTenant(this.tenant);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
}

Test.prototype.AddTenant = function () {
    const response = this.api.AddTenant(this.tenant, this.tenant);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
}

Test.prototype.UploadWasmFile = function (fileId, path) {
    const response = this.api.UploadWasmFile(this.tenant, fileId, path);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
}

Test.prototype.UploadSchemaFile = function (fileId, path) {
    const response = this.api.UploadSchemaFile(this.tenant, fileId, path);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

}

Test.prototype.AddPackageWithFile = function (path) {
    const response = this.api.AddPackageFile(path);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
}

Test.prototype.AddPackageFileForJobWithTemplate = function (id, name, idQ, nameQ, path) {
    const job = this.api.GetPackageObj(path);
    job.package.ID = id;
    job.package.name = name;
    job.package.queues[0].ID = idQ;
    job.package.queues[0].name = nameQ;
    const response = this.api.InvokeAddPackage(job);
    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });
}

Test.prototype.Close = function () {
    this.api.Close();
}