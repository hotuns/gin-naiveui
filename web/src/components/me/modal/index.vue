<template>
  <n-modal
    v-model:show="show"
    :style="{ width: modalOptions.width, ...modalOptions.modalStyle }"
    :preset="undefined"
    size="huge"
    :bordered="false"
  >
    <n-card
      :title="modalOptions.title"
      :style="modalOptions.contentStyle"
      :closable="modalOptions.closable"
      @close="close()"
    >
      <slot></slot>

      <!-- 底部按钮 -->
      <template #footer>
        <slot v-if="$slots.footer" name="footer" />
        <footer v-else-if="modalOptions.showFooter" class="flex justify-end">
          <n-button v-if="modalOptions.showCancel" @click="handleCancel()">
            {{ modalOptions.cancelText }}
          </n-button>
          <n-button
            v-if="modalOptions.showOk"
            type="primary"
            :loading="modalOptions.okLoading"
            class="ml-20"
            @click="handleOk()"
          >
            {{ modalOptions.okText }}
          </n-button>
        </footer>
      </template>
    </n-card>
  </n-modal>
</template>

<script setup>
const props = defineProps({
  width: {
    type: String,
    default: '800px',
  },
  title: {
    type: String,
    default: '',
  },
  closable: {
    type: Boolean,
    default: true,
  },
  cancelText: {
    type: String,
    default: '取消',
  },
  okText: {
    type: String,
    default: '确定',
  },
  showFooter: {
    type: Boolean,
    default: true,
  },
  showCancel: {
    type: Boolean,
    default: true,
  },
  showOk: {
    type: Boolean,
    default: true,
  },
  modalStyle: {
    type: Object,
    default: () => {},
  },
  contentStyle: {
    type: Object,
    default: () => {},
  },
  onOk: {
    type: Function,
    default: () => {},
  },
  onCancel: {
    type: Function,
    default: () => {},
  },
})
// 声明一个show变量，用于控制模态框的显示与隐藏
const show = ref(false)
// 声明一个modalOptions变量，用于存储模态框的配置信息
const modalOptions = ref({})

// 打开模态框
function open(options = {}) {
  // 将props和options合并赋值给modalOptions
  modalOptions.value = { ...props, ...options }

  // 将show的值设置为true
  show.value = true
}

// 定义一个close函数，用于关闭模态框
function close() {
  show.value = false
}

// 定义一个handleOk函数，用于处理模态框确定操作
async function handleOk(data) {
  // 如果modalOptions中没有onOk函数，则直接关闭模态框
  if (typeof modalOptions.value.onOk !== 'function') {
    return close()
  }
  try {
    // 调用onOk函数，传入data参数
    const res = await modalOptions.value.onOk(data)
    // 如果onOk函数的返回值不为false，则关闭模态框
    res !== false && close()
  } catch (error) {
    okLoading.value = false
    console.error(error)
  }
}

// 定义一个handleCancel函数，用于处理模态框取消操作
async function handleCancel(data) {
  // 如果modalOptions中没有onCancel函数，则直接关闭模态框
  if (typeof modalOptions.value.onCancel !== 'function') {
    return close()
  }
  try {
    // 调用onCancel函数，传入data参数
    const res = await modalOptions.value.onCancel(data)

    // 如果onCancel函数的返回值不为false，则关闭模态框
    res !== false && close()
  } catch (error) {
    okLoading.value = false
    console.error(error)
  }
}

const okLoading = computed({
  get() {
    return !!modalOptions.value?.okLoading
  },
  set(v) {
    if (modalOptions.value) {
      modalOptions.value.okLoading = v
    }
  },
})

// 定义一个defineExpose函数，用于暴露open、close、handleOk、handleCancel函数
defineExpose({
  open,
  close,
  handleOk,
  handleCancel,
  okLoading,
  options: modalOptions,
})
</script>
