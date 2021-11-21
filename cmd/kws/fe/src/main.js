// 主应用
import { createApp } from 'vue'
import App from './App.vue'
let app = createApp(App)

// 路由规则
import Home from '@/views/Home.vue'
const routes = [
  { path: '/', name: 'Home', component: Home, meta: { keepAlive: true } }
]
import { createRouter, createWebHashHistory } from 'vue-router'
const router = createRouter({
  history: createWebHashHistory(),
  routes: routes
})
app.use(router)

// font awesome 图标
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faHeart } from '@fortawesome/free-solid-svg-icons'
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons'
import { faPlus } from '@fortawesome/free-solid-svg-icons'
import { faMinus } from '@fortawesome/free-solid-svg-icons'
import { faHeart as faHeartRegular } from '@fortawesome/free-regular-svg-icons'
library.add(faHeart)
library.add(faArrowLeft)
library.add(faHeartRegular)
library.add(faPlus)
library.add(faMinus)
app.component('font-awesome-icon', FontAwesomeIcon)

// 饿了么组件
import 'element-plus/lib/theme-chalk/display.css'
import 'element-plus/lib/theme-chalk/base.css'
import {
  ElDialog,
  ElDivider,
  ElInput,
  ElIcon,
  ElLink,
  ElContainer,
  ElUpload,
  ElForm,
  ElCol,
  ElRow,
  ElMain,
  ElHeader,
  ElFormItem,
  ElButton,
  ElSelect,
  ElSwitch,
  ElOption,
  ElTable,
  ElTableColumn,
  ElRadio,
  ElAside,
  ElPopconfirm,
  ElImage,
  ElDatePicker,
  ElConfigProvider,
  ElTooltip,
  ElBacktop,
  ElCheckbox,
  ElTag
} from 'element-plus'
app.use(ElDialog)
app.use(ElDivider)
app.use(ElInput)
app.use(ElIcon)
app.use(ElLink)
app.use(ElContainer)
app.use(ElUpload)
app.use(ElForm)
app.use(ElRow)
app.use(ElCol)
app.use(ElMain)
app.use(ElHeader)
app.use(ElFormItem)
app.use(ElButton)
app.use(ElSelect)
app.use(ElSwitch)
app.use(ElOption)
app.use(ElTable)
app.use(ElTableColumn)
app.use(ElRadio)
app.use(ElAside)
app.use(ElPopconfirm)
app.use(ElImage)
app.use(ElDatePicker)
app.use(ElConfigProvider)
app.use(ElTooltip)
app.use(ElBacktop)
app.use(ElCheckbox)
app.use(ElTag)

// 加载主应用
app.mount('#app')
