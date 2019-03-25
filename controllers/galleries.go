package controllers

import "sw.com/FirstWebWithGO/views"

type Galleries struct {
	Gallery *views.View
}

func NewGalleries() *Galleries {
	return &Galleries{
		Gallery: views.NewView("bootstrap", "users/new"),
	}
}
