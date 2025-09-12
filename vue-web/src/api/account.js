import { get, post, postForm } from './http'

export function profile() { return get('/account/profile') }
export function packages() { return get('/account/packages') }
export function usage(range) { return get('/account/usage', range) }
export function login(phone, password) { return postForm('/auth/login', { phone, password }) }
export function changePassword(new_password) { return postForm('/auth/change_password', { new_password }) }
export function logout(token) { return post(`/auth/logout?token=${encodeURIComponent(token)}`) }



