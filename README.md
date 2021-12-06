# kubecost_exporter

This is the PoC version of kubecost-exporter. The kubecost is able to export out-of-cluster metrics to a prometheus, but there is no detailed information about labels and properties of exported objects. This exporter uses [Asets API](https://github.com/kubecost/docs/blob/master/assets.md) to export data from the cost-model to prometheus. 
Also I was told by the kubecost support team that they will deprecate the out-of-cluster costs export to a prometheus, so I decided to write my own exporter.
<img width="635" alt="support_reply" src="https://user-images.githubusercontent.com/3328394/144870424-48b54bec-7ccc-4e7c-909e-964d53a785ce.png">


As for now it doesn't support any filtering and can only export all assets data provided by the kubecost Assets API. 
