import Home from './views/Home.vue'
import NotFound from './views/NotFound.vue'
import HistoryTTS from './views/HistoryTTS.vue'
import Account from './views/Account.vue'

/** @type {import('vue-router').RouterOptions['routes']} */
export let routes = [
  { path: '/', component: Home, meta: { title: '首页' } },
  { path: '/history/tts', component: HistoryTTS, meta: { title: '合成历史' } },
  { path: '/account', component: Account, meta: { title: '账号' } },
  { path: '/:path(.*)', component: NotFound },
]
