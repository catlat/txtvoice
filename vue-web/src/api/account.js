import { get, post, postForm } from './http'

export function profile() { return get('/account/profile') }
export function packages() { return get('/account/packages') }
export function usage(range) { return get('/account/usage', range) }
export function loginSimple(identity) { return postForm('/auth/login_simple', { identity }) }
export function logout(token) { return post(`/auth/logout?token=${encodeURIComponent(token)}`) }



