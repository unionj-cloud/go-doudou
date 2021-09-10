package codegen

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

var statefulsetTmpl = `apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{.SvcName}}-statefulset
spec:
  selector:
    matchLabels:
      app: {{.SvcName}}
  serviceName: {{.SvcName}}-svc-headless
  replicas: 1
  template:
    metadata:
      labels:
        app: {{.SvcName}}
    spec:
      terminationGracePeriodSeconds: 10
      containers:
        - name: {{.SvcName}}
          image: {{.Image}}
          imagePullPolicy: Always
          ports:
            - name: http-port
              containerPort: 6060
              protocol: TCP
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  name: {{.SvcName}}-svc-headless
spec:
  selector:
    app: {{.SvcName}}
  ports:
    - protocol: TCP
      port: 6060
      targetPort: 6060
  clusterIP: None`

func GenK8sStatefulset(dir string, svcname, image string) {
	var (
		f   *os.File
		tpl *template.Template
	)
	file := filepath.Join(dir, svcname+"_statefulset.yaml")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		if f, err = os.Create(file); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("statefulset.tmpl").Parse(statefulsetTmpl); err != nil {
			panic(err)
		}
		if err = tpl.Execute(f, struct {
			SvcName string
			Image   string
		}{
			SvcName: svcname,
			Image:   image,
		}); err != nil {
			panic(err)
		}
	} else {
		logrus.Warnf("image version will be modified in file %s", file)
		err = ioutil.WriteFile(file, modifyVersion(file, image), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
