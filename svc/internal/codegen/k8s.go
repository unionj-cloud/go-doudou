package codegen

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var k8sTmpl = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.SvcName}}-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.SvcName}}
  template:
    metadata:
      labels:
        app: {{.SvcName}}
    spec:
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
  name: {{.SvcName}}-service
spec:
  type: LoadBalancer
  externalTrafficPolicy: Cluster
  selector:
    app: {{.SvcName}}
  ports:
    - protocol: TCP
      port: 6060
      targetPort: 6060`

func GenK8s(dir string, svcname, image string) {
	var (
		f   *os.File
		tpl *template.Template
	)
	k8sfile := filepath.Join(dir, svcname+"_k8s.yaml")
	if fi, err := os.Stat(k8sfile); os.IsNotExist(err) {
		if f, err = os.Create(k8sfile); err != nil {
			panic(err)
		}
		defer f.Close()

		if tpl, err = template.New("k8sfile.tmpl").Parse(k8sTmpl); err != nil {
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
		logrus.Warnf("file %s will be overwrite", k8sfile)
		err = ioutil.WriteFile(k8sfile, modifyVersion(k8sfile, image), fi.Mode()&os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func modifyVersion(yfile string, image string) []byte {
	var (
		f                             *os.File
		err                           error
		raw, jdeployment, ddeployment []byte
		deployment                    string
	)
	if f, err = os.Open(yfile); err != nil {
		panic(err)
	}
	defer f.Close()
	raw, err = ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	blocks := strings.Split(string(raw), "---")
	deployment = blocks[0]

	jdeployment, err = yaml.YAMLToJSON([]byte(deployment))
	if err != nil {
		panic(err)
	}
	c, err := gabs.ParseJSON(jdeployment)
	if err != nil {
		panic(err)
	}
	c.Set(image, gabs.DotPathToSlice("spec.template.spec.containers.0.image")...)
	if err != nil {
		panic(err)
	}
	ddeployment, _ = yaml.JSONToYAML(c.Bytes())
	blocks[0] = string(ddeployment)
	return []byte(strings.Join(blocks, "---"))
}
