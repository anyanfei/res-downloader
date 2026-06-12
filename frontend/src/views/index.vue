<template>
  <div class="h-full flex flex-col px-5 pt-5 overflow-y-auto [&::-webkit-scrollbar]:hidden">
    <Transition name="slide-up">
      <div v-if="showBatchProgress" class="fixed right-5 bottom-5 w-80 rounded-xl shadow-lg border border-gray-200 bg-white dark:bg-gray-800 dark:border-gray-700 z-50" style="--wails-draggable:no-drag">
        <div class="flex items-center justify-between px-4 py-3 border-b border-gray-100 dark:border-gray-700">
          <span class="text-sm font-medium">{{ t('index.batch_progress') }}</span>
          <NButton size="tiny" quaternary @click="showBatchProgress = false">
            <template #icon>
              <NIcon><CloseOutline /></NIcon>
            </template>
          </NButton>
        </div>
        <div class="px-4 py-3">
          <div class="mb-2 text-xs text-gray-500">{{ batchCompleted + batchError }} / {{ batchTotal }}</div>
          <NProgress
            type="line"
            :percentage="batchCompleted + batchError === batchTotal ? 100 : Math.round((batchCompleted + batchError) / Math.max(batchTotal, 1) * 100)"
            :status="batchError > 0 ? 'warning' : 'success'"
            :indicator-placement="'inside'"
            :height="22"
            :border-radius="4"
            processing
          />
          <div class="flex justify-between mt-1 text-xs text-gray-500">
            <span>{{ t('index.done') }}: {{ batchCompleted }}</span>
            <span v-if="batchError > 0">{{ t('index.error') }}: {{ batchError }}</span>
          </div>
        </div>
      </div>
    </Transition>
    <div class="pb-2 z-40" id="header">
      <NSpace>
        <NButton v-if="isProxy" secondary type="primary" @click.stop="close" style="--wails-draggable:no-drag">
          <span class="inline-block w-1.5 h-1.5 bg-red-600 rounded-full mr-1 animate-pulse"></span>
          {{ t("index.close_grab") }}{{ data.length > 0 ? `&nbsp;${t('index.total_resources', {count: data.length})}` : '' }}
        </NButton>
        <NButton v-else tertiary type="tertiary" @click.stop="open" style="--wails-draggable:no-drag">
          {{ t("index.open_grab") }}{{ data.length > 0 ? `&nbsp;${t('index.total_resources', {count: data.length})}` : '' }}
        </NButton>
        <NSelect style="min-width: 100px;--wails-draggable:no-drag" :placeholder="t('index.grab_type')" v-model:value="resourcesType" multiple clearable
                 :max-tag-count="3" :options="classify"></NSelect>
        <NSelect
            style="min-width: 100px;--wails-draggable:no-drag"
            :placeholder="t('index.pin_type')"
            v-model:value="pinClassify"
            :options="classify.slice(1)"
            clearable
        />
        <NButtonGroup style="--wails-draggable:no-drag">

          <NButton v-if="rememberChoice" tertiary type="error" @click.stop="clear" style="--wails-draggable:no-drag">
            <template #icon>
              <n-icon>
                <TrashOutline/>
              </n-icon>
            </template>
            {{ t("index.clear_list") }}
          </NButton>
          <n-popconfirm
              v-else
              @positive-click="()=>{rememberChoice=rememberChoiceTmp;clear()}"
              :show-icon="false"
          >
            <template #trigger>
              <NButton tertiary type="error" style="--wails-draggable:no-drag">
                <template #icon>
                  <n-icon>
                    <TrashOutline/>
                  </n-icon>
                </template>
                {{ t("index.clear_list") }}
              </NButton>
            </template>
            <div>
              <div class="flex flex-row items-center text-red-700 my-2 text-base">
                <n-icon>
                  <TrashOutline/>
                </n-icon>
                <p class="ml-1">{{ t("index.clear_list_tip") }}</p>
              </div>
              <NCheckbox
                  v-model:checked="rememberChoiceTmp"
              >
                <span class="text-gray-400">{{ t('index.remember_clear_choice') }}</span>
              </NCheckbox>
            </div>
          </n-popconfirm>

          <NButton tertiary type="primary" @click.stop="batchDown">
            <template #icon>
              <n-icon>
                <DownloadOutline/>
              </n-icon>
            </template>
            {{ t('index.batch_download') }}
          </NButton>
          <NButton tertiary type="info">
            <NPopover placement="bottom" trigger="hover">
              <template #trigger>
                <NIcon size="18" class="">
                  <Apps/>
                </NIcon>
              </template>
              <div class="flex flex-col">
                <NButton tertiary type="error" @click.stop="batchCancel" class="my-1">
                  <template #icon>
                    <n-icon>
                      <CloseOutline/>
                    </n-icon>
                  </template>
                  {{ t('index.cancel_down') }}
                </NButton>
                <NButton tertiary type="warning" @click.stop="batchExport()" class="my-1">
                  <template #icon>
                    <n-icon>
                      <ArrowRedoCircleOutline/>
                    </n-icon>
                  </template>
                  {{ t('index.batch_export') }}
                </NButton>
                <NButton tertiary type="info" @click.stop="showImport=true" class="my-1">
                  <template #icon>
                    <n-icon>
                      <ServerOutline/>
                    </n-icon>
                  </template>
                  {{ t('index.batch_import') }}
                </NButton>
                <NButton tertiary type="primary" @click.stop="batchExport('url')" class="my-1">
                  <template #icon>
                    <n-icon>
                      <ArrowRedoCircleOutline/>
                    </n-icon>
                  </template>
                  {{ t('index.export_url') }}
                </NButton>
              </div>
            </NPopover>
          </NButton>
        </NButtonGroup>
      </NSpace>
    </div>
    <div class="flex-1">
      <NDataTable
          :columns="columns"
          :data="filteredData"
          :bordered="false"
          :max-height="tableHeight"
          :row-key="rowKey"
          :virtual-scroll="true"
          :header-height="48"
          :height-for-row="()=> 48"
          :checked-row-keys="checkedRowKeysValue"
          @update:checked-row-keys="handleCheck"
          @update:filters="updateFilters"
          style="--wails-draggable:no-drag"
      />
    </div>
    <div class="flex justify-center items-center text-blue-400" id="bottom">
      <span class="cursor-pointer px-2 py-1" @click="BrowserOpenURL(certUrl)">{{ t('footer.cert_download') }}</span>
      <span class="cursor-pointer px-2 py-1" @click="BrowserOpenURL('https://github.com/putyy/res-downloader')">{{ t('footer.source_code') }}</span>
      <span class="cursor-pointer px-2 py-1" @click="BrowserOpenURL('https://github.com/putyy/res-downloader/issues')">{{ t('footer.help') }}</span>
      <span class="cursor-pointer px-2 py-1" @click="BrowserOpenURL('https://github.com/putyy/res-downloader/releases')">{{ t('footer.update_log') }}</span>
    </div>
    <Preview v-model:showModal="showPreviewRow" :previewRow="previewRow"/>
    <ShowLoading :loadingText="loadingText" :isLoading="loading"/>
    <ImportJson v-model:showModal="showImport" @submit="handleImport"/>
    <Password v-model:showModal="showPassword" @submit="handlePassword"/>
  </div>
</template>

<script lang="ts" setup>
import {NButton, NIcon, NImage, NInput, NSpace, NTooltip, NPopover, NGradientText, NProgress} from "naive-ui"
import {computed, h, onMounted, ref, watch} from "vue"
import type {appType} from "@/types/app"
import type {DataTableRowKey, ImageRenderToolbarProps, DataTableFilterState, DataTableBaseColumn} from "naive-ui"
import Preview from "@/components/Preview.vue"
import ShowLoading from "@/components/ShowLoading.vue"
// @ts-ignore
import {getDecryptionArray} from '@/assets/js/decrypt.js'
import {useIndexStore} from "@/stores"
import appApi from "@/api/app"
import Action from "@/components/Action.vue"
import ActionDesc from "@/components/ActionDesc.vue"
import ImportJson from "@/components/ImportJson.vue"
import {useEventStore} from "@/stores/event"
import {BrowserOpenURL, ClipboardSetText} from "../../wailsjs/runtime"
import Password from "@/components/Password.vue"
import ShowOrEdit from "@/components/ShowOrEdit.vue"
import {useI18n} from 'vue-i18n'
import {
  DownloadOutline,
  ArrowRedoCircleOutline,
  ServerOutline,
  SearchOutline,
  Apps,
  TrashOutline, CloseOutline
} from "@vicons/ionicons5"
import {useDialog} from 'naive-ui'
import * as bind from "../../wailsjs/go/core/Bind"
import {Quit} from "../../wailsjs/runtime"
import {DialogOptions} from "naive-ui/es/dialog/src/DialogProvider"
import {formatSize} from "@/func"

const {t} = useI18n()
const eventStore = useEventStore()
const dialog = useDialog()
const isProxy = computed(() => {
  return store.isProxy
})
const certUrl = computed(() => {
  return store.baseUrl + "/api/cert"
})
const data = ref<any[]>([])
const filterClassify = ref<string[]>([])
const pinClassify = ref<string | null>(null)
const filteredData = computed(() => {
  let result = data.value

  if (filterClassify.value.length > 0) {
    result = result.filter(item => filterClassify.value.includes(item.Classify))
  }

  if (descriptionSearchValue.value) {
    result = result.filter(item => item.Description?.toLowerCase().includes(descriptionSearchValue.value.toLowerCase()))
  }

  if (urlSearchValue.value) {
    result = result.filter(item => item.Url?.toLowerCase().includes(urlSearchValue.value.toLowerCase()))
  }

  if (pinClassify.value) {
    const pinnedItems = result.filter(item => item.Classify === pinClassify.value)
    const otherItems = result.filter(item => item.Classify !== pinClassify.value)
    result = [...pinnedItems, ...otherItems]
  }

  return result
})

const store = useIndexStore()
const tableHeight = ref(800)
const resourcesType = ref<string[]>(["all"])

const classifyAlias: { [key: string]: any } = {
  image: computed(() => t("index.image")),
  audio: computed(() => t("index.audio")),
  video: computed(() => t("index.video")),
  m3u8: computed(() => t("index.m3u8")),
  live: computed(() => t("index.live")),
  xls: computed(() => t("index.xls")),
  doc: computed(() => t("index.doc")),
  pdf: computed(() => t("index.pdf")),
  stream: computed(() => t("index.stream")),
  font: computed(() => t("index.font"))
}

const dwStatus = computed<any>(() => {
  return {
    ready: t("index.ready"),
    pending: t("index.pending"),
    running: t("index.running"),
    error: t("index.error"),
    done: t("index.done"),
    handle: t("index.handle")
  }
})

const maxConcurrentDownloads = computed(() => {
  return store.globalConfig.DownNumber
})

const classify = ref([
  {
    value: "all",
    label: computed(() => t("index.all")),
  },
])

const descriptionSearchValue = ref("")
const urlSearchValue = ref("")
const rememberChoice = ref(false)
const rememberChoiceTmp = ref(false)

const columns = ref<any[]>([
  {
    type: "selection",
  },
  {
    title: () => {
      if (checkedRowKeysValue.value.length > 0) {
        return h(NGradientText, {type: "success"}, t("index.choice") + `(${checkedRowKeysValue.value.length})`)
      }
      return h('div', {class: 'flex items-center'}, [
        t('index.domain'),
        h(NPopover, {
          style: "--wails-draggable:no-drag",
          trigger: 'click',
          placement: 'bottom',
          showArrow: true,
        }, {
          trigger: () => h(NIcon, {
            size: "18",
            class: `ml-1 cursor-pointer ${urlSearchValue.value ? "text-green-600": "text-gray-500"}`,
            onClick: (e: MouseEvent) => e.stopPropagation()
          }, h(SearchOutline)),
          default: () => h('div', {class: 'p-2 w-64'}, [
            h(NInput, {
              value: urlSearchValue.value,
              'onUpdate:value': (val: string) => urlSearchValue.value = val,
              placeholder: t('index.search_description'),
              clearable: true
            }, {
              prefix: () => h(NIcon, {component: SearchOutline})
            })
          ])
        })
      ])
    },
    key: "Domain",
    width: 90,
    render: (row: appType.MediaInfo) => {
      return h(NTooltip, {
        trigger: 'hover',
        placement: 'top'
      }, {
        trigger: () => h('span', {
          class: 'cursor-default'
        }, row.Domain),
        default: () => row.Url
      })
    }
  },
  {
    title: computed(() => t("index.type")),
    key: "Classify",
    width: 80,
    filterOptions: computed(() => Array.from(classify.value).slice(1)),
    filterMultiple: true,
    filter: (value: string, row: appType.MediaInfo) => {
      return !!~row.Classify.indexOf(String(value))
    },
    render: (row: appType.MediaInfo) => {
      const item = classify.value.find(item => item.value === row.Classify)
      const label = item ? item.label : row.Classify
      if (pinClassify.value && row.Classify === pinClassify.value) {
        return h('span', {class: 'pin-flashing'}, label)
      }
      return label
    }
  },
  {
    title: computed(() => t("index.preview")),
    key: "Url",
    width: 80,
    render: (row: appType.MediaInfo) => {
      if (row.Classify === "image") {
        return h("div", {
          style: "width: 100%;max-height:80px;overflow:hidden;"
        }, h(NImage, {
          objectFit: "contain",
          lazy: true,
          "render-toolbar": renderToolbar,
          src: row.Url,
          onError: () => {
            const idx = data.value.findIndex(item => item.Id === row.Id)
            if (idx !== -1) {
              data.value.splice(idx, 1)
              cacheData()
            }
          },
          onLoad: (e: Event) => {
            const img = e.target as HTMLImageElement
            if (img && img.naturalWidth <= 1 && img.naturalHeight <= 1) {
              const idx = data.value.findIndex(item => item.Id === row.Id)
              if (idx !== -1) {
                data.value.splice(idx, 1)
                cacheData()
              }
            }
          }
        }))
      }
      return [
        h(
            NButton,
            {
              strong: true,
              tertiary: true,
              type: "info",
              size: "small",
              style: {
                margin: "2px"
              },
              onClick: () => {
                if (row.Classify === "audio" || row.Classify === "video" || row.Classify === "m3u8" || row.Classify === "live") {
                  previewRow.value = row
                  showPreviewRow.value = true
                }
              }
            },
            {
              default: () => {
                if (row.Classify === "audio" || row.Classify === "video" || row.Classify === "m3u8" || row.Classify === "live") {
                  return t("index.preview")
                }
                return t("index.preview_tip")
              }
            }
        ),
      ]
    }
  },
  {
    title: computed(() => t("index.status")),
    key: "Status",
    width: 80,
    render: (row: appType.MediaInfo, index: number) => {
      let status = "info"
      if (row.Status === "done" || row.Status === "running") {
        status = "success"
      } else if (row.Status === "pending") {
        status = "warning"
      }

      return h(
          NButton,
          {
            tertiary: true,
            type: status as any,
            size: "small",
            style: {
              margin: "2px"
            },
            onClick: () => {
              if (row.SavePath && row.Status === "done") {
                appApi.openFolder({filePath: row.SavePath})
              } else if (row.Status === "ready") {
                download(row, index)
              }
            }
          },
          {
            default: () => {
              return row.Status === "running" ? row.SavePath : dwStatus.value[row.Status as keyof typeof dwStatus]
            }
          }
      )
    }
  },
  {
    title: () => h('div', {class: 'flex items-center'}, [
      t('index.description'),
      h(NPopover, {
        style: "--wails-draggable:no-drag",
        trigger: 'click',
        placement: 'bottom',
        showArrow: true,
      }, {
        trigger: () => h(NIcon, {
          size: "18",
          class: `ml-1 cursor-pointer ${descriptionSearchValue.value ? "text-green-600": "text-gray-500"}`,
          onClick: (e: MouseEvent) => e.stopPropagation()
        }, h(SearchOutline)),
        default: () => h('div', {class: 'p-2 w-64'}, [
          h(NInput, {
            value: descriptionSearchValue.value,
            'onUpdate:value': (val: string) => descriptionSearchValue.value = val,
            placeholder: t('index.search_description'),
            clearable: true
          }, {
            prefix: () => h(NIcon, {component: SearchOutline})
          })
        ])
      })
    ]),
    key: "Description",
    width: 150,
    render: (row: appType.MediaInfo, index: number) => {
      return h(ShowOrEdit, {
        value: row.Description,
        onUpdateValue(v: string) {
          data.value[index].Description = v
          cacheData()
        }
      })
    }
  },
  {
    title: computed(() => t("index.resource_size")),
    key: "Size",
    width: 120,
    sorter: (row1: appType.MediaInfo, row2: appType.MediaInfo) => row1.Size - row2.Size,
    render(row: appType.MediaInfo, index: number) {
      return formatSize(row.Size)
    }
  },
  {
    title: computed(() => t("index.save_path")),
    key: "SavePath",
    render(row: appType.MediaInfo, index: number) {
      return h("a",
          {
            href: "javascript:;",
            class: "ellipsis-2",
            style: {
              color: "#5a95d0"
            },
            onClick: () => {
              if (row.SavePath && row.Status === "done") {
                appApi.openFolder({filePath: row.SavePath})
              }
            }
          },
          row.Status === "running" ? "" : row.SavePath
      )
    }
  },
  {
    key: "actions",
    width: 130,
    render(row: appType.MediaInfo, index: number) {
      return h(Action, {key: index, row: row, index: index, onAction: dataAction})
    },
    title() {
      return h(ActionDesc)
    }
  }
])

const checkedRowKeysValue = ref<DataTableRowKey[]>([])
const showPreviewRow = ref(false)
const previewRow = ref<appType.MediaInfo>()
const loading = ref(false)
const loadingText = ref("")
const showImport = ref(false)
const showPassword = ref(false)
const downloadQueue = ref<appType.MediaInfo[]>([])
let activeDownloads = 0
let isOpenProxy = false
let isInstall = false
const batchTotal = ref(0)
const batchCompleted = ref(0)
const batchError = ref(0)
const showBatchProgress = ref(false)

onMounted(() => {
  try {
    window.addEventListener("resize", () => {
      resetTableHeight()
    })
    loading.value = true
    handleInstall().then((is: boolean) => {
      isInstall = true
      loading.value = false
    })

    checkLoading()
    watch(showPassword, () => {
      if (!showPassword.value) {
        checkLoading()
      }
    })
  } catch (e) {
    window.$message?.error(JSON.stringify(e), {duration: 5000})
  }

  buildClassify()

  const temp = localStorage.getItem("resources-type")
  if (temp) {
    resourcesType.value = JSON.parse(temp).res
  } else {
    appApi.setType(resourcesType.value)
  }

  const cache = localStorage.getItem("resources-data")
  if (cache) {
    data.value = JSON.parse(cache)
  }

  const choiceCache = localStorage.getItem("remember-clear-choice")
  if (choiceCache === "1") {
    rememberChoice.value = true
  }

  const pinClassifyCache = localStorage.getItem("pin-classify")
  if (pinClassifyCache) {
    pinClassify.value = pinClassifyCache
  }

  watch(rememberChoice, (n, o) => {
    if (rememberChoice.value) {
      localStorage.setItem("remember-clear-choice", "1")
    } else {
      localStorage.removeItem("remember-clear-choice")
    }
  })

  watch(pinClassify, (n, o) => {
    if (pinClassify.value) {
      localStorage.setItem("pin-classify", pinClassify.value)
    } else {
      localStorage.removeItem("pin-classify")
    }
  })

  resetTableHeight()

  eventStore.addHandle({
    type: "newResources",
    event: (res: appType.MediaInfo) => {
      if (store.globalConfig.InsertTail) {
        data.value.push(res)
      } else {
        data.value.unshift(res)
      }
      cacheData()
    }
  })

  eventStore.addHandle({
    type: "downloadProgress",
    event: (res: { Id: string, SavePath: string, Status: string, Message: string }) => {
      switch (res.Status) {
        case "running":
          updateItem(res.Id, item => {
            if (item.Status !== 'done') {
              item.SavePath = res.Message
              item.Status = 'running'
            }
          })
          break
        case "done":
          updateItem(res.Id, item => {
            item.SavePath = res.SavePath
            item.Status = 'done'
          })
          activeDownloads = Math.max(0, activeDownloads - 1)
          if (showBatchProgress.value) {
            batchCompleted.value++
            checkBatchComplete()
          }
          cacheData()
          checkQueue()
          break
        case "error":
          updateItem(res.Id, item => {
            item.SavePath = res.Message
            item.Status = 'error'
          })
          activeDownloads = Math.max(0, activeDownloads - 1)
          if (showBatchProgress.value) {
            batchError.value++
            checkBatchComplete()
          }
          cacheData()
          checkQueue()
          break
      }
    }
  })
})

watch(() => {
  return store.globalConfig.MimeMap
}, () => {
  buildClassify()
})

watch(resourcesType, (n, o) => {
  localStorage.setItem("resources-type", JSON.stringify({res: resourcesType.value}))
  appApi.setType(resourcesType.value)
})

const updateItem = (id: string, updater: (item: any) => void) => {
  const item = data.value.find(i => i.Id === id)
  if (item) updater(item)
}

function cacheData() {
  localStorage.setItem("resources-data", JSON.stringify(data.value))
}

const resetTableHeight = () => {
  try {
    const headerHeight = document.getElementById("header")?.offsetHeight || 0
    const bottomHeight = document.getElementById("bottom")?.offsetHeight || 0
    // @ts-ignore
    const theadHeight = document.getElementsByClassName("n-data-table-thead")[0]?.offsetHeight || 0
    const height = document.documentElement.clientHeight || window.innerHeight
    tableHeight.value = height - headerHeight - bottomHeight - theadHeight - 20
  } catch (e) {
    console.log(e)
  }
}

const buildClassify = () => {
  const mimeMap = store.globalConfig.MimeMap ?? {}
  const seen = new Set()
  const types = Object.values(mimeMap)
      .filter(({Type}) => {
        if (seen.has(Type)) return false
        seen.add(Type)
        return true
      })
      .map(({Type}) => ({
        value: Type,
        label: classifyAlias[Type] ?? Type,
      }))
      .sort((a, b) => {
        if (a.value === 'video') return -1
        if (b.value === 'video') return 1
        return 0
      })
  classify.value = [
    {value: "all", label: computed(() => t("index.all"))},
    ...types,
  ]
}

const dataAction = (row: appType.MediaInfo, index: number, type: string) => {
  switch (type) {
    case "down":
      download(row, index)
      break
    case "cancel":
      if (row.Status === "pending") {
        const queueIndex = downloadQueue.value.findIndex(item => item.Id === row.Id)
        if (queueIndex !== -1) {
          downloadQueue.value.splice(queueIndex, 1)
        }
        updateItem(row.Id, item => {
          item.Status = 'ready'
          item.SavePath = ''
        })
        cacheData()
      } else if (row.Status === "running") {
        appApi.cancel({id: row.Id}).then((res) => {
          updateItem(row.Id, item => {
            item.Status = 'ready'
            item.SavePath = ''
          })
          if (activeDownloads > 0) {
            activeDownloads--
          }
          cacheData()
          checkQueue()
          if (res.code === 0) {
            window?.$message?.error(res.message)
            return
          }
        })
      }
      break
    case "copy":
      ClipboardSetText(row.Url).then((is: boolean) => {
        if (is) {
          window?.$message?.success(t("common.copy_success"))
        } else {
          window?.$message?.error(t("common.copy_fail"))
        }
      })
      break
    case "json":
      ClipboardSetText(encodeURIComponent(JSON.stringify(row))).then((is: boolean) => {
        if (is) {
          window?.$message?.success(t("common.copy_success"))
        } else {
          window?.$message?.error(t("common.copy_fail"))
        }
      })
      break
    case "open":
      BrowserOpenURL(row.Url)
      break
    case "decode":
      decodeWxFile(row, index)
      break
    case "delete":
      if (row.Status === "pending" || row.Status === "running") {
        window?.$message?.error(t("index.delete_tip"))
        return
      }
      appApi.delete({sign: [row.UrlSign]}).then(() => {
        const idx = data.value.findIndex(item => item.Id === row.Id)
        if (idx !== -1) {
          data.value.splice(idx, 1)
        }
        cacheData()
      })
      break
  }
}

const renderToolbar = ({nodes}: ImageRenderToolbarProps) => {
  return [
    nodes.rotateCounterclockwise,
    nodes.rotateClockwise,
    nodes.resizeToOriginalSize,
    nodes.zoomOut,
    nodes.zoomIn,
    nodes.close
  ]
}

const rowKey = (row: appType.MediaInfo) => {
  return row.Id
}

const handleCheck = (rowKeys: DataTableRowKey[]) => {
  checkedRowKeysValue.value = rowKeys
}

const updateFilters = (filters: DataTableFilterState, initiatorColumn: DataTableBaseColumn) => {
  filterClassify.value = filters.Classify as string[]
}

const batchDown = async () => {
  if (checkedRowKeysValue.value.length <= 0) {
    window?.$message?.error(t("index.use_data"))
    return
  }

  if (!store.globalConfig.SaveDirectory) {
    window?.$message?.error(t("index.save_path_empty"))
    return
  }

  const items = data.value.filter(item =>
    checkedRowKeysValue.value.includes(item.Id) && item.Classify !== 'live' && item.Classify !== 'm3u8'
  )
  batchTotal.value = items.length
  batchCompleted.value = 0
  batchError.value = 0
  showBatchProgress.value = true

  items.forEach((item, index) => {
    download(item, index)
  })

  checkedRowKeysValue.value = []
}

const playCompleteSound = () => {
  try {
    const ctx = new AudioContext();
    // 五声音阶，像风铃一样
    const notes = [523, 587, 659, 784, 880, 1047, 1175, 1319, 1568, 1760];
    // 随机挑几个不同高度的音
    const picked = [0, 2, 4, 7, 9].map(i => notes[i]);
    // 叠加轻微延迟，模拟风吹动
    const delays = [0, 80, 150, 250, 330];
    picked.forEach((freq, i) => {
      const t = ctx.currentTime + delays[i] / 1000 + Math.random() * 0.04;
      const harmFreqs = [freq, freq * 2.76, freq * 5.4];
      harmFreqs.forEach((hf, hi) => {
        const osc = ctx.createOscillator();
        const gain = ctx.createGain();
        osc.type = hi === 0 ? 'triangle' : 'sine';
        osc.frequency.value = hf;
        const vol = [0.3, 0.12, 0.04][hi];
        gain.gain.setValueAtTime(0, t);
        gain.gain.linearRampToValueAtTime(vol, t + 0.01);
        gain.gain.exponentialRampToValueAtTime(0.001, t + 1.2);
        osc.connect(gain);
        gain.connect(ctx.destination);
        osc.start(t);
        osc.stop(t + 1.2);
      });
    });
    setTimeout(() => ctx.close(), 1800);
  } catch (_) {}
}

const checkBatchComplete = () => {
  if (showBatchProgress.value && batchCompleted.value + batchError.value === batchTotal.value) {
    playCompleteSound()
  }
}

const batchCancel = async () => {
  if (checkedRowKeysValue.value.length <= 0) {
    window?.$message?.error(t("index.use_data"))
    return
  }
  loading.value = true
  const cancelTasks: Promise<any>[] = []
  data.value.forEach((item, index) => {
    if (!checkedRowKeysValue.value.includes(item.Id)) {
      return
    }

    if (item.Status === "pending") {
      const queueIndex = downloadQueue.value.findIndex(qItem => qItem.Id === item.Id)
      if (queueIndex !== -1) {
        downloadQueue.value.splice(queueIndex, 1)
      }
      item.Status = 'ready'
      item.SavePath = ''
      return
    }

    if (item.Status === "running") {
      activeDownloads = Math.max(0, activeDownloads - 1)
      cancelTasks.push(appApi.cancel({id: item.Id}).then(() => {
        item.Status = 'ready'
        item.SavePath = ''
        checkQueue()
      }))
    }
  })
  await Promise.allSettled(cancelTasks)
  loading.value = false
  checkedRowKeysValue.value = []
  cacheData()
}

const batchExport = (type?: string) => {
  if (checkedRowKeysValue.value.length <= 0) {
    window?.$message?.error(t("index.use_data"))
    return
  }

  if (!store.globalConfig.SaveDirectory) {
    window?.$message?.error(t("index.save_path_empty"))
    return
  }

  loadingText.value = t("common.loading")
  loading.value = true

  let jsonData = data.value.filter(item => checkedRowKeysValue.value.includes(item.Id))

  if (type === "url") {
    jsonData = jsonData.map(item => item.Url)
  } else {
    jsonData = jsonData.map(item => encodeURIComponent(JSON.stringify(item)))
  }

  appApi.batchExport({content: jsonData.join("\n")}).then((res: appType.Res) => {
    loading.value = false
    if (res.code === 0) {
      window?.$message?.error(res.message)
      return
    }
    window?.$message?.success(t("index.import_success"))
    window?.$message?.info(t("index.save_path") + "：" + res.data?.file_name, {
      duration: 5000
    })
  })
}

const uint8ArrayToBase64 = (bytes: any) => {
  return window.btoa(Array.from(bytes, (byte: any) => String.fromCharCode(byte)).join(''))
}

const download = (row: appType.MediaInfo, index: number) => {
  if (!store.globalConfig.SaveDirectory) {
    window?.$message?.error(t("index.save_path_empty"))
    return
  }

  if (data.value.some(item => item.Id === row.Id && item.Status === "running")) {
    return
  }

  if (downloadQueue.value.some(item => item.Id === row.Id || item.Url === row.Url)) {
    return
  }

  if (activeDownloads >= maxConcurrentDownloads.value) {
    row.Status = "pending"
    downloadQueue.value.push(row)
    return
  }

  if (row.Status === "done" || row.Status === "error") {
    row.Status = "ready"
    row.SavePath = ""
  }

  startDownload(row, index)
}

const startDownload = (row: appType.MediaInfo, index: number) => {
  activeDownloads++

  const decodeStr = row.DecodeKey
      ? uint8ArrayToBase64(getDecryptionArray(row.DecodeKey))
      : ""

  appApi.download({...row, decodeStr}).then((res: appType.Res) => {
    if (res.code === 0) {
      window?.$message?.error(res.message)
    }
  })
}

const checkQueue = () => {
  if (downloadQueue.value.length > 0 && activeDownloads < maxConcurrentDownloads.value) {
    const nextItem = downloadQueue.value.shift()
    if (nextItem) {
      const index = data.value.findIndex(item => item.Id === nextItem.Id)
      if (index !== -1) {
        startDownload(nextItem, index)
      }
    }
  }
}

const open = () => {
  isOpenProxy = true
  store.openProxy().then((res: appType.Res) => {
    if (res.code === 1) {
      return
    }

    if (["darwin", "linux"].includes(store.envInfo.platform)) {
      showPassword.value = true
    } else {
      window.$message?.error(res.message)
    }
  })
}

const close = () => {
  store.unsetProxy()
}

const clear = async () => {
  const newData = [] as any[]
  const signs: string[] = []
  if (checkedRowKeysValue.value.length > 0) {
    data.value.forEach((item, index) => {
      if (checkedRowKeysValue.value.includes(item.Id) && item.Status !== "pending" && item.Status !== "running") {
        signs.push(item.UrlSign)
      } else {
        newData.push(item)
      }
    })
    checkedRowKeysValue.value = []
  } else {
    data.value.forEach((item, index) => {
      if (item.Status === "pending" || item.Status === "running") {
        newData.push(item)
      } else {
        signs.push(item.UrlSign)
      }
    })
  }
  await appApi.delete({sign: signs})
  data.value = newData
  cacheData()
}

const decodeWxFile = (row: appType.MediaInfo, index: number) => {
  if (!row.DecodeKey) {
    window?.$message?.error(t("index.video_decode_no"))
    return
  }
  appApi.openFileDialog().then((res: appType.Res) => {
    if (res.code === 0) {
      window?.$message?.error(res.message)
      return
    }
    if (res.data.file) {
      loadingText.value = t("index.video_decode_loading")
      loading.value = true
      appApi.wxFileDecode({
        ...row,
        filename: res.data.file,
        decodeStr: uint8ArrayToBase64(getDecryptionArray(row.DecodeKey))
      }).then((res: appType.Res) => {
        loading.value = false
        if (res.code === 0) {
          window?.$message?.error(res.message)
          return
        }
        data.value[index].SavePath = res.data.save_path
        data.value[index].Status = "done"
        cacheData()
        window?.$message?.success(t("index.video_decode_success"))
      })
    }
  })
}

const handleImport = (content: string) => {
  if (!content) {
    window?.$message?.error(t("view.import_empty"))
    return
  }
  let newItems = [] as any[]
  content.split("\n").forEach((line, index) => {
    try {
      let res = JSON.parse(decodeURIComponent(line))
      if (res && res?.Id) {
        res.Id = res.Id + Math.floor(Math.random() * 100000)
        res.SavePath = ""
        res.Status = "ready"
        newItems.push(res)
      }
    } catch (e) {
      console.log(e)
    }
  })
  if (newItems.length > 0) {
    data.value = [...newItems, ...data.value]
    cacheData()
  }
  showImport.value = false
}

const handlePassword = async (password: string, isCache: boolean) => {
  const res = await appApi.setSystemPassword({password, isCache})
  if (res.code === 0) {
    window.$message?.error(res.message)
    return
  }

  if (isOpenProxy) {
    showPassword.value = false
    store.openProxy()
    return
  }

  handleInstall().then((is: boolean) => {
    if (is) {
      showPassword.value = false
    }
  })
}

const handleInstall = async () => {
  isOpenProxy = false
  const res = await appApi.install()
  if (res.code === 1) {
    store.globalConfig.AutoProxy && store.openProxy()
    return true
  }

  window.$message?.error(res.message, {duration: 5000})

  if (store.envInfo.platform === "windows" && res.message.includes("Access is denied")) {
    window.$message?.error(t("index.win_install_tip"))
  } else if (["darwin", "linux"].includes(store.envInfo.platform)) {
    showPassword.value = true
  }
  return false
}

const checkLoading = () => {
  setTimeout(() => {
    if (loading.value && !isInstall && !showPassword.value) {
      dialog.warning({
        title: t("index.start_err_tip"),
        content: t("index.start_err_content"),
        positiveText: t("index.start_err_positiveText"),
        negativeText: t("index.start_err_negativeText"),
        draggable: false,
        closeOnEsc: false,
        closable: false,
        maskClosable: false,
        onPositiveClick: () => {
          bind.ResetApp()
        },
        onNegativeClick: () => {
          Quit()
        }
      } as DialogOptions)
    }
  }, 6000)
}
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(120%);
  opacity: 0;
}
</style>