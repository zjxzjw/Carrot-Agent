import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import Backend from 'i18next-http-backend'

import zhCN from './locales/zh-CN.json'
import enUS from './locales/en-US.json'

const resources = {
  'zh-CN': {
    translation: zhCN
  },
  'en-US': {
    translation: enUS
  }
}

i18n
  .use(Backend)
  .use(initReactI18next)
  .init({
    resources,
    lng: 'zh-CN',
    fallbackLng: 'en-US',
    interpolation: {
      escapeValue: false
    },
    detection: {
      order: ['localStorage', 'navigator']
    }
  })

export default i18n