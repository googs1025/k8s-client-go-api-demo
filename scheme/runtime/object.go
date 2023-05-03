package runtime

type Object interface {
	GetObjectKind(g GroupVersionKind) (ObjectKind, error)
}

type ObjectKind interface {
	SetGroupVersionKind(kind GroupVersionKind)
	GroupVersionKind() GroupVersionKind
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}
