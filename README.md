# workflew

events
workers


Package definition:

- Name
- MerchantId
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
merchantId: soytriguerov1
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