package app

//func TestServerApp_Run(t *testing.T) {
//	router := chi.NewRouter()
//	logger := loggers.NewLogger()
//	type fields struct {
//		router *chi.Mux
//		cfg    config.ServerConfig
//		logger *loggers.Logger
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{
//			name: "ok",
//			fields: fields{
//				router: router,
//				logger: logger,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			a := &ServerApp{
//				router: tt.fields.router,
//				cfg:    tt.fields.cfg,
//				logger: tt.fields.logger,
//			}
//			go a.Run()
//		})
//	}
//}
