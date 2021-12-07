# kubecost_exporter

This is the PoC version of kubecost-exporter. The kubecost is able to export out-of-cluster metrics to a prometheus, but there is no detailed information about labels and properties of exported objects. This exporter uses [Asets API](https://github.com/kubecost/docs/blob/master/assets.md) to export data from the cost-model to prometheus. 
Also I was told by the kubecost support team that they will deprecate the out-of-cluster costs export to a prometheus, so I decided to write my own exporter.

<img width="621" alt="Screenshot 2021-12-06 at 17 10 40" src="https://user-images.githubusercontent.com/3328394/144870829-9496cd3a-bff0-4af6-965f-e7c3beb06931.png">


As for now it doesn't support any filtering and can only export all assets data provided by the kubecost Assets API. 

---
### TODO list
- Write tests!!
- Add the ability to handle user defined filters. 
- Add basic auth support for the KubeCost client api.
- Refactor some parts of code marked with _TODO_ labels. (and maybe something else)
- Add something to this list :)
