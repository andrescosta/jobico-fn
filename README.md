# github.com/andrescosta/jobico

events
workers


Package definition:

- Name
- TenantId
- Queues[]:
    - Queue
        - Name

- Executors[]
    - Executor
        - Name
        - Package
            - Type (WASM function, types: ...)
        - SupportedEvents[]
            - Event
                - Name

- Events[]
    - Event
        - Name
        - Schema
    - SupplierQueues[]
        - Queue
            - Name

YAML:

name: mypackage
tenantId: soytriguerov1
queues:
  -
    name: queue1
    queueId: aaa

events:
  -
    name:
    schema:
        schemaId: ll
        name:
        schemaRef:
    supplierQueue:
        name: queue1


executors:
  -
    name: myexs
    package:
        packageId:
        name:
        packageRef:
        type: wasi
    supportedevents:
      -
        name: evt1

notification:
    - queue
    - notificator

### YAML Definition

name: mypackage
jobpackageid: mypackage
tenantid: m1
queues:
  - queueid: queue1
    name: queue1
  - queueid: queue2
    name: queue2

events:
  - name: ev1
    eventid: ev1
    datatype: 0
    schema:
      schemaid: sche1
      name: sche1
      schemaref: schema.json
    supplierqueueid: queue1
    runtimeid: runtime1
  - name: ev2
    eventid: ev2
    datatype: 0
    schema:
      schemaid: sche1
      name: sche1
      schemaref: schema.json
    supplierqueueid: queue2
    runtimeid: runtime1
        

runtimes:
  - runtimeid: runtime1
    name: greet.wasm
    moduleref: greet.wasm
    mainfuncname: event
    type: 0