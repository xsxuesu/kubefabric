package kubeutils
import (
	"errors"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func (k *KubeClient)CreatePv(pvName,labelname,server string,path string )(*apiv1.PersistentVolume,error){

	pv := k.Client.CoreV1().PersistentVolumes()
	pvment := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: pvName,
			Labels: map[string]string{
				"type":"local",
				"name":labelname,
			},
		},
		Spec: apiv1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy:apiv1.PersistentVolumeReclaimRecycle,
			AccessModes:                   []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany},
			//MountOptions:                  apiv1.MountOptionAnnotation,
			Capacity: apiv1.ResourceList{
				apiv1.ResourceName(apiv1.ResourceStorage):resource.MustParse("10Gi"),
			},
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				NFS: &apiv1.NFSVolumeSource{
					Server:   server,
					Path:     path,
					ReadOnly: false,
				},
			},
		},
	}
	persistent,err := pv.Create(pvment)
	if err != nil {
		return nil, err
	}
	return persistent, nil
}

func (k *KubeClient)WatchPv(pvName string)(int,error){
	pv := k.Client.CoreV1().PersistentVolumes()
	winter,err := pv.Watch(metav1.ListOptions{})
	if err != nil {
		return 0,err
	}
	select {
	case wr := <- winter.ResultChan():
		switch wr.Type {
		case watch.Added:
			fmt.Println(wr.Object)
			return 1,nil
		case watch.Error:
			fmt.Println(wr.Object)
			return 0,errors.New("create pv err")
		case watch.Deleted:
			fmt.Println(wr.Object)
			return -1,nil
		case watch.Modified:
			fmt.Println(wr.Object)
			return 1,nil
		}
	}
	return 0,nil
}

func (k *KubeClient)UpdatePv(pvName string)(*apiv1.PersistentVolume,error){
	pv := k.Client.CoreV1().PersistentVolumes()
	pvconfig, err := pv.Get(pvName,metav1.GetOptions{})
	if err != nil {
		return nil,err
	}
	pvcon , err := pv.Update(pvconfig)
	if err != nil {
		return nil,err
	}
	return pvcon,nil
}

func (k *KubeClient)GetPvList()([]apiv1.PersistentVolume, error) {
	pv := k.Client.CoreV1().PersistentVolumes()
	pvList,err := pv.List(metav1.ListOptions{})
	if err != nil {
		return nil,err
	}
	for _,pv := range pvList.Items {
		fmt.Printf("persistent volumn : %s namespace:%s \n",pv.Name,pv.Namespace)
	}
	return pvList.Items,nil
}

func (k *KubeClient)DeletePv(pvName string)error{
	pv := k.Client.CoreV1().PersistentVolumes()
	err := pv.Delete(pvName,&metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (k *KubeClient)CreatePVC(namespace,pvcName string,selectorlabel map[string]string)(*apiv1.PersistentVolumeClaim,error){
	//storage := "standard"
	volumeMode := apiv1.PersistentVolumeFilesystem
	pvcInter := k.Client.CoreV1().PersistentVolumeClaims(namespace)
	pvcment  := &apiv1.PersistentVolumeClaim{
		ObjectMeta:metav1.ObjectMeta{
			Name: pvcName,
		},
		Spec:apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteMany},
			VolumeMode : &volumeMode,
			//StorageClassName:&storage,
			Resources:apiv1.ResourceRequirements{
				Requests:apiv1.ResourceList{
					apiv1.ResourceName(apiv1.ResourceStorage):resource.MustParse("10Gi"),
				},
			},
			//VolumeName:pvcName,
			Selector: &metav1.LabelSelector{
				MatchLabels:selectorlabel,
			},
		},
	}

	persistent,err := pvcInter.Create(pvcment)
	if err != nil {
		return nil,err
	}
	return persistent,nil
}

func (k *KubeClient)WatchPvc(namespace string)(int,error){
	pvc := k.Client.CoreV1().PersistentVolumeClaims(namespace)
	winter,err := pvc.Watch(metav1.ListOptions{})
	if err != nil {
		return 0,err
	}
	select {
	case wr := <- winter.ResultChan():
		switch wr.Type {
		case watch.Added:
			fmt.Println(wr.Object)
			return 1,nil
		case watch.Error:
			fmt.Println(wr.Object)
			return 0,errors.New("create pvc err")
		case watch.Deleted:
			fmt.Println(wr.Object)
			return -1,nil
		case watch.Modified:
			fmt.Println(wr.Object)
			return 1,nil
		}
	}
	return 0,nil
}

func (k *KubeClient)UpdatePvc(namespace string,pvcname string)(*apiv1.PersistentVolumeClaim,error){
	pvcMent := k.Client.CoreV1().PersistentVolumeClaims(namespace)
	pvc,err := pvcMent.Get(pvcname,metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	pvc.Labels = map[string]string{}

	pvc2,err := pvcMent.Update(pvc)
	if err != nil {
		return nil , err
	}
	return pvc2, nil
}

func (k *KubeClient)GetPvcList(namespace string)([]apiv1.PersistentVolumeClaim,error){
	pvcList,err := k.Client.CoreV1().PersistentVolumeClaims(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil,err
	}
	for _,pvc := range pvcList.Items{
		fmt.Printf("persistent volumn claim : %s namespace:%s \n",pvc.Name,pvc.Namespace)
	}
	return pvcList.Items,nil
}

func (k *KubeClient)DelectPvc(namespace,pvcName string)error{
	pvcMent := k.Client.CoreV1().PersistentVolumeClaims(namespace)
	err := pvcMent.Delete(pvcName,&metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
