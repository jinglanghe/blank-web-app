package service

import (
	"context"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func initInformer() {
	ctx := context.TODO()
	informerFactory := informers.NewSharedInformerFactory(clientset, 0)

	initNodeInformer(ctx, informerFactory)
	initCMInformer(ctx, informerFactory)
}

func initNodeInformer(ctx context.Context, informerFactory informers.SharedInformerFactory) {
	nodeInformer := informerFactory.Core().V1().Nodes().Informer()
	nodeInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc:    nodeInformerHandlerAddFunc,
		UpdateFunc: nodeInformerHandlerUpdateFunc,
		DeleteFunc: nodeInformerHandlerDeleteFunc,
	})

	go informerFactory.Start(ctx.Done())

	cache.WaitForCacheSync(ctx.Done(), nodeInformer.HasSynced)
}

func nodeInformerHandlerAddFunc(obj interface{}) {
	node := obj.(*v1.Node)
	logging.Debug().Msgf("node add: %v", node.Name)
	RefreshData()
}

func nodeInformerHandlerUpdateFunc(oldObj, newObj interface{}) {
	oldNode := oldObj.(*v1.Node)
	newNode := newObj.(*v1.Node)
	logging.Debug().Msgf("node update: oldNode.name:%v, newNode.name:%v", oldNode.Name, newNode.Name)
	RefreshData()
}

func nodeInformerHandlerDeleteFunc(obj interface{}) {
	node := obj.(*v1.Node)
	logging.Debug().Msgf("node delete: %v", node.Name)
	RefreshData()
}

func initCMInformer(ctx context.Context, informerFactory informers.SharedInformerFactory) {
	configMapInformer := informerFactory.Core().V1().ConfigMaps().Informer()
	configMapInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: configMapInformerHandlerAddFunc,
	})

	go informerFactory.Start(ctx.Done())

	cache.WaitForCacheSync(ctx.Done(), configMapInformer.HasSynced)
}

func configMapInformerHandlerAddFunc(obj interface{}) {
	cm := obj.(*v1.ConfigMap)
	if cm.Name != PromConfigmap {
		return
	}

	logging.Debug().Msgf("configMap add: %v", cm.Name)
}
