The [New Relic Go Agent](https://github.com/newrelic/go-agent) provides a nice middleware for the stdlib handler signature. 
The following is an adaptation of that middleware for Gin.

```golang
const (
	// NewRelicTxnKey is the key used to retrieve the NewRelic Transaction from the context
	NewRelicTxnKey = "NewRelicTxnKey"
)

// NewRelicMonitoring is a middleware that starts a newrelic transaction, stores it in the context, then calls the next handler
func NewRelicMonitoring(app newrelic.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		txn := app.StartTransaction(ctx.Request.URL.Path, ctx.Writer, ctx.Request)
		defer txn.End()
		ctx.Set(NewRelicTxnKey, txn)
		ctx.Next()
	}
}
```
and in `main.go` or equivalent...
```golang
router := gin.Default()
cfg := newrelic.NewConfig(os.Getenv("APP_NAME"), os.Getenv("NEW_RELIC_API_KEY"))
app, err := newrelic.NewApplication(cfg)
if err != nil {
		log.Printf("failed to make new_relic app: %v", err)
} else {
		router.Use(adapters.NewRelicMonitoring(app))
}
 ```
