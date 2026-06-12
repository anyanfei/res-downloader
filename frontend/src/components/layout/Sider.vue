<template>
  <div class="flex pb-2 flex-col h-full min-w-[80px] border-r border-slate-100 dark:border-slate-900">
    <Screen v-if="envInfo.platform!=='darwin'"></Screen>
    <div class="w-full flex flex-row items-center justify-center pt-5" :class="envInfo.platform==='darwin' ? 'pt-8' : 'pt-2'">
      <div class="relative flex items-center justify-center cursor-pointer" @click="showAppInfo = true">
        <img class="w-12 h-12 rounded-full transition-transform duration-300 hover:scale-105 dark" src="@/assets/image/logo.png" alt="res-downloader logo"/>
      </div>
    </div>
    <main class="flex-1 flex-grow-1 mb-5 overflow-auto flex flex-col pt-1 items-center h-full" v-if="is">
      <NScrollbar :size="1">
        <NLayout has-sider>
          <NLayoutSider
              :bordered="false"
              show-trigger
              collapse-mode="width"
              :on-after-enter="() => { showAppName = true }"
              :on-after-leave="() => { showAppName = false }"
              :collapsed-width="70"
              :collapsed="collapsed"
              :width="envInfo.platform==='linux' ? 160 : 140"
              :native-scrollbar="false"
              :inverted="inverted"
              :on-update:collapsed="collapsedChange"
              class="bg-inherit"
          >
            <NMenu
                :inverted="inverted"
                :collapsed-width="70"
                :collapsed-icon-size="22"
                :options="menuOptions"
                :value="menuValue"
                @update:value="handleUpdateValue"
            />
          </NLayoutSider>
        </NLayout>
        <NLayoutFooter position="absolute" :inverted="inverted" class="bg-inherit">
          <NMenu
              :inverted="inverted"
              :collapsed-width="70"
              :collapsed-icon-size="22"
              :options="footerOptions"
              :value="menuValue"
              @update:value="handleFooterUpdate"
          />
        </NLayoutFooter>
      </NScrollbar>
    </main>
  </div>
  <Footer v-model:showModal="showAppInfo"/>
</template>

<script lang="ts" setup>
import {MenuOption} from "naive-ui"
import {NIcon} from "naive-ui"
import {computed, h, onMounted, ref, watch} from "vue"
import {useRoute, useRouter} from "vue-router"
import {
  CloudOutline,
  SettingsOutline,
  HelpCircleOutline,
  MoonOutline,
  SunnyOutline,
  LanguageSharp,
  ShieldCheckmarkOutline,
  SyncOutline
} from "@vicons/ionicons5"
import {useIndexStore} from "@/stores"
import Footer from "@/components/Footer.vue"
import Screen from "@/components/Screen.vue"
import {useI18n} from "vue-i18n"
import appApi from "@/api/app"

const {t} = useI18n()
const route = useRoute()
const router = useRouter()
const inverted = ref(false)
const collapsed = ref(false)
const showAppName = ref(false)
const showAppInfo = ref(false)
const menuValue = ref(route.fullPath.substring(1))
const store = useIndexStore()
const is = ref(false)
const certInstalling = ref(false)

const envInfo = store.envInfo

const globalConfig = computed(() => {
  return store.globalConfig
})

const theme = computed(() => {
  return store.globalConfig.Theme === "darkTheme" ? renderIcon(SunnyOutline) : renderIcon(MoonOutline)
})

const certOptions = computed(() => {
  if (certInstalling.value) {
    return {
      label: t("footer.cert_installing"),
      icon: () => h(NIcon, {class: 'spin-icon'}, {default: () => h(SyncOutline)}),
    }
  }
  return {
    label: t("footer.cert_download"),
    icon: renderIcon(ShieldCheckmarkOutline),
  }
})

watch(() => route.path, (newPath, oldPath) => {
  menuValue.value = route.fullPath.substring(1)
})

onMounted(()=>{
  const collapsedCache = localStorage.getItem("collapsed");
  if (collapsedCache) {
    collapsed.value = JSON.parse(collapsedCache).collapsed
  }
  is.value = true
})

const renderIcon = (icon: any) => {
  return () => h(NIcon, null, {default: () => h(icon)})
}

const menuOptions = ref([
  {
    label: computed(() => t("menu.index")),
    key: 'index',
    icon: renderIcon(CloudOutline),
  },
  {
    label: computed(() => t("menu.setting")),
    key: 'setting',
    icon: renderIcon(SettingsOutline),
  },
])

const footerOptions = computed(() => [
  {
    label: certOptions.value.label,
    key: 'cert',
    icon: certOptions.value.icon,
  },
  {
    label: t("menu.locale"),
    key: 'locale',
    icon: renderIcon(LanguageSharp),
  },
  {
    label: t("menu.theme"),
    key: 'theme',
    icon: theme.value,
  },
  {
    label: t("menu.about"),
    key: 'about',
    icon: renderIcon(HelpCircleOutline),
  },
])

const handleUpdateValue = (key: string, item?: MenuOption) => {
  menuValue.value = key
  return router.push({path: "/" + key})
}

const handleCertInstall = async () => {
  if (certInstalling.value) return
  certInstalling.value = true
  try {
    const checkRes = await appApi.certCheck()
    if (checkRes.code === 1 && checkRes.data?.installed) {
      window?.$message?.success(t('footer.cert_installed'))
      return
    }
    const res = await appApi.install()
    if (res.code === 1) {
      window?.$message?.success(t('footer.cert_install_success'))
    } else {
      window?.$message?.error(res.message, {duration: 5000})
      if (store.envInfo.platform === "windows" && res.message.includes("Access is denied")) {
        window?.$message?.error(t('index.win_install_tip'))
      }
    }
  } catch (e: any) {
    window?.$message?.error(String(e))
  } finally {
    certInstalling.value = false
  }
}

const handleFooterUpdate = (key: string, item?: MenuOption) => {
  if (key === "about") {
    showAppInfo.value = true
    return
  }

  if (key === "cert") {
    handleCertInstall()
    return
  }

  if (key === "theme") {
    if (globalConfig.value.Theme === "darkTheme") {
      store.setConfig({Theme: "lightTheme"})
      return
    }
    store.setConfig({Theme: "darkTheme"})
    return
  }

  if (key === "locale") {
    if (globalConfig.value.Locale === "zh") {
      store.setConfig({Locale: "en"})
      return
    }
    store.setConfig({Locale: "zh"})
    return
  }

  menuValue.value = key
  return router.push({path: "/" + key})
}

const collapsedChange = (value: boolean)=>{
  collapsed.value = value
  localStorage.setItem("collapsed", JSON.stringify({collapsed: value}))
}
</script>
<style scoped>
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.spin-icon {
  animation: spin 1s linear infinite;
}
</style>
