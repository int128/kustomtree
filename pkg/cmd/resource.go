package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ResourceMetadata struct {
	Name string `yaml:"name"`
}

func processResource(name string, basedir string) error {
	fullpath := filepath.Join(basedir, name)
	s, err := os.Stat(fullpath)
	if err != nil {
		return fmt.Errorf("could not stat: %w", err)
	}
	if s.IsDir() {
		return nil
	}

	f, err := os.Open(fullpath)
	if err != nil {
		return fmt.Errorf("could not open the file: %w", err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	for {
		var n yaml.Node
		if err := d.Decode(&n); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("could not decode the file: %w", err)
		}
		if err := processResourceNode(name, n); err != nil {
			return err
		}
	}
	return nil
}

func processResourceNode(name string, n yaml.Node) error {
	var v struct {
		Kind     string           `yaml:"kind"`
		Metadata ResourceMetadata `yaml:"metadata"`
	}
	if err := n.Decode(&v); err != nil {
		return fmt.Errorf("could not decode the node: %w", err)
	}

	switch v.Kind {
	case "Deployment":
		desiredName := fmt.Sprintf("deployment/%s.yaml", v.Metadata.Name)
		if name != desiredName {
			log.Printf("TODO: you need to extract %s: %s -> %s", v.Kind, name, desiredName)
		}
	case "CronJob":
		desiredName := fmt.Sprintf("cronjob/%s.yaml", v.Metadata.Name)
		if name != desiredName {
			log.Printf("TODO: you need to extract %s: %s -> %s", v.Kind, name, desiredName)
		}
	case "Service":
		desiredName := fmt.Sprintf("service/%s.yaml", v.Metadata.Name)
		if name != desiredName {
			log.Printf("TODO: you need to extract %s: %s -> %s", v.Kind, name, desiredName)
		}
	case "HorizontalPodAutoscaler":
		desiredName := fmt.Sprintf("hpa/%s.yaml", v.Metadata.Name)
		if name != desiredName {
			log.Printf("TODO: you need to extract %s: %s -> %s", v.Kind, name, desiredName)
		}
	case "PodDisruptionBudget":
		desiredName := fmt.Sprintf("pdb/%s.yaml", v.Metadata.Name)
		if name != desiredName {
			log.Printf("TODO: you need to extract %s: %s -> %s", v.Kind, name, desiredName)
		}
	case "ConfigMap":
		desiredName := fmt.Sprintf("configmap/%s.yaml", v.Metadata.Name)
		if name != desiredName {
			log.Printf("TODO: you need to extract %s: %s -> %s", v.Kind, name, desiredName)
		}
	default:
		log.Printf("unknown kind %s", v.Kind)
	}
	return nil
}

//func encodeResourceNode(n yaml.Node) error {
//	e := yaml.NewEncoder(os.Stdout)
//	e.SetIndent(2)
//	if err := e.Encode(&n); err != nil {
//		return fmt.Errorf("could not encode the node: %w", err)
//	}
//	return nil
//}
