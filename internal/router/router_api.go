// Package router
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-09-04 16:35
package router

func setApiRouter(r *resource) {
	//// helper
	//helpers := r.mux.Group("/helper")
	//{
	//
	//}
	//
	//// admin
	//
	//// 需要签名验证，无需登录验证，无需 RBAC 权限验证
	//login := r.mux.Group("/api", r.interceptors.CheckSignature())
	//{
	//}
	//
	//// 需要签名验证、登录验证，无需 RBAC 权限验证
	//notRBAC := r.mux.Group("/api", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckSignature())
	//{
	//
	//}
	//
	//// 需要签名验证、登录验证、RBAC 权限验证
	//api := r.mux.Group("/api", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckSignature(), r.interceptors.CheckRBAC())
	//{
	//}
}
