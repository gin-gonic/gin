package yaml

import "context"

type (
	ctxMergeKey  struct{}
	ctxAnchorKey struct{}
)

func withMerge(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxMergeKey{}, true)
}

func isMerge(ctx context.Context) bool {
	v, ok := ctx.Value(ctxMergeKey{}).(bool)
	if !ok {
		return false
	}
	return v
}

func withAnchor(ctx context.Context, name string) context.Context {
	anchorMap := getAnchorMap(ctx)
	if anchorMap == nil {
		anchorMap = make(map[string]struct{})
	}
	anchorMap[name] = struct{}{}
	return context.WithValue(ctx, ctxAnchorKey{}, anchorMap)
}

func getAnchorMap(ctx context.Context) map[string]struct{} {
	v, ok := ctx.Value(ctxAnchorKey{}).(map[string]struct{})
	if !ok {
		return nil
	}
	return v
}
