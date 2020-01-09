package util

import (
	"encoding/json"
	"errors"

	hawtiov1alpha1 "github.com/hawtio/hawtio-operator/pkg/apis/hawtio/v1alpha1"

	consolev1 "github.com/openshift/api/console/v1"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetHawtconfig reads the console configuration from the config map
func GetHawtconfig(configMap *corev1.ConfigMap) (*hawtiov1alpha1.Hawtconfig, error) {
	var hawtconfig *hawtiov1alpha1.Hawtconfig

	hawtconfigJSON, ok := configMap.Data["hawtconfig.json"]
	if !ok {
		return hawtconfig, errors.New("did not find hawtconfig.json in ConfigMap")
	}
	err := json.Unmarshal([]byte(hawtconfigJSON), &hawtconfig)
	if err != nil {
		return hawtconfig, err
	}

	return hawtconfig, nil
}

// NewConsoleLink creates a ConsoleLink instance
func NewConsoleLink(name string, route *routev1.Route, hawtconfig *hawtiov1alpha1.Hawtconfig) *consolev1.ConsoleLink {
	consoleLink := &consolev1.ConsoleLink{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"app": "hawtio"},
		},
		Spec: consolev1.ConsoleLinkSpec{
			Location:        consolev1.ApplicationMenu,
			ApplicationMenu: &consolev1.ApplicationMenuSpec{},
		},
	}

	UpdateLink(consoleLink, route, hawtconfig)

	return consoleLink
}

// UpdateLink updates the console link properties
func UpdateLink(consoleLink *consolev1.ConsoleLink, route *routev1.Route, hawtconfig *hawtiov1alpha1.Hawtconfig) {
	consoleLink.Spec.Link.Text = hawtconfig.Branding.ConsoleLink.Text
	consoleLink.Spec.Link.Href = "https://" + route.Spec.Host
	consoleLink.Spec.ApplicationMenu.Section = hawtconfig.Branding.ConsoleLink.Section
	consoleLink.Spec.ApplicationMenu.ImageURL = "https://" + route.Spec.Host + hawtconfig.Branding.ConsoleLink.ImageRelativePath
}