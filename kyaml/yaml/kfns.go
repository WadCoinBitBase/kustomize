// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/kyaml/errors"
)

// AnnotationClearer removes an annotation at metadata.annotations.
// Returns nil if the annotation or field does not exist.
type AnnotationClearer struct {
	Kind string `yaml:"kind,omitempty"`
	Key  string `yaml:"key,omitempty"`
}

func (c AnnotationClearer) Filter(rn *RNode) (*RNode, error) {
	return rn.Pipe(
		PathGetter{Path: []string{"metadata", "annotations"}},
		FieldClearer{Name: c.Key})
}

func ClearAnnotation(key string) AnnotationClearer {
	return AnnotationClearer{Key: key}
}

// ClearEmptyAnnotations clears the keys, annotations
// and metadata if they are empty/null
func ClearEmptyAnnotations(rn *RNode) error {
	_, err := rn.Pipe(Lookup("metadata"), FieldClearer{
		Name: "annotations", IfEmpty: true})
	if err != nil {
		return errors.Wrap(err)
	}
	_, err = rn.Pipe(FieldClearer{Name: "metadata", IfEmpty: true})
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

// AnnotationSetter sets an annotation at metadata.annotations.
// Creates metadata.annotations if does not exist.
type AnnotationSetter struct {
	Kind  string `yaml:"kind,omitempty"`
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

func (s AnnotationSetter) Filter(rn *RNode) (*RNode, error) {
	// some tools get confused about the type if annotations are not quoted
	v := NewScalarRNode(s.Value)
	v.YNode().Tag = StringTag
	v.YNode().Style = yaml.SingleQuotedStyle

	if err := ClearEmptyAnnotations(rn); err != nil {
		return nil, err
	}

	return rn.Pipe(
		PathGetter{Path: []string{"metadata", "annotations"}, Create: yaml.MappingNode},
		FieldSetter{Name: s.Key, Value: v})
}

func SetAnnotation(key, value string) AnnotationSetter {
	return AnnotationSetter{Key: key, Value: value}
}

// AnnotationGetter gets an annotation at metadata.annotations.
// Returns nil if metadata.annotations does not exist.
type AnnotationGetter struct {
	Kind  string `yaml:"kind,omitempty"`
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

// AnnotationGetter returns the annotation value.
// Returns "", nil if the annotation does not exist.
func (g AnnotationGetter) Filter(rn *RNode) (*RNode, error) {
	v, err := rn.Pipe(PathGetter{Path: []string{"metadata", "annotations", g.Key}})
	if v == nil || err != nil {
		return v, err
	}
	if g.Value == "" || v.value.Value == g.Value {
		return v, err
	}
	return nil, err
}

func GetAnnotation(key string) AnnotationGetter {
	return AnnotationGetter{Key: key}
}

// LabelSetter sets a label at metadata.labels.
// Creates metadata.labels if does not exist.
type LabelSetter struct {
	Kind  string `yaml:"kind,omitempty"`
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}

func (s LabelSetter) Filter(rn *RNode) (*RNode, error) {
	// some tools get confused about the type if labels are not quoted
	v := NewScalarRNode(s.Value)
	v.YNode().Tag = StringTag
	v.YNode().Style = yaml.SingleQuotedStyle
	return rn.Pipe(
		PathGetter{Path: []string{"metadata", "labels"}, Create: yaml.MappingNode},
		FieldSetter{Name: s.Key, Value: v})
}

func SetLabel(key, value string) LabelSetter {
	return LabelSetter{Key: key, Value: value}
}
