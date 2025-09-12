<template>
  <!-- 遮罩层 -->
  <div v-if="show" class="absolute inset-0 bg-white/70 backdrop-blur-xs rounded-xl z-10 flex items-center justify-center min-h-[200px]">
    <div class="text-center max-w-sm mx-auto px-6 py-8 bg-white/85 rounded-xl shadow-lg backdrop-blur-sm">
        <!-- 主要动画图标 -->
        <div class="mb-4">
          <div v-if="type === 'tts'" class="mb-6">
            <!-- 音波动画 -->
            <div class="flex items-end justify-center space-x-1.5 h-12">
              <div v-for="i in 5" :key="i" 
                   class="w-1.5 bg-gradient-to-t from-gray-400 to-gray-600 rounded-full animate-wave"
                   :class="getBarHeight(i)"
                   :style="{ animationDelay: `${i * 0.15}s` }">
              </div>
            </div>
          </div>

          <!-- 视频获取动画 -->
          <div v-else-if="type === 'video'" class="mb-6">
            <div class="flex items-center justify-center">
              <div class="relative">
                <!-- 外圈旋转 -->
                <div class="w-12 h-12 border-2 border-gray-300 border-t-gray-600 rounded-full animate-spin"></div>
                <!-- 中心视频图标 -->
                <div class="absolute inset-0 flex items-center justify-center">
                  <svg class="w-6 h-6 text-gray-600" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M2 6a2 2 0 012-2h6l2 2h6a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14 9a2 2 0 11-4 0 2 2 0 014 0z"/>
                  </svg>
                </div>
              </div>
            </div>
          </div>

          <!-- 字幕提取动画 -->
          <div v-else-if="type === 'subtitle'" class="mb-6">
            <div class="flex items-center justify-center relative">
              <!-- 多层圆环动画 -->
              <div class="relative w-16 h-16">
                <!-- 外圈：慢速旋转 -->
                <div class="absolute inset-0 w-16 h-16 border-2 border-blue-200 border-t-blue-500 rounded-full animate-spin" style="animation-duration: 3s;"></div>
                <!-- 中圈：中速旋转 -->
                <div class="absolute inset-2 w-12 h-12 border-2 border-indigo-200 border-r-indigo-500 rounded-full animate-spin" style="animation-duration: 2s; animation-direction: reverse;"></div>
                <!-- 内圈：快速旋转 -->
                <div class="absolute inset-4 w-8 h-8 border-2 border-purple-200 border-b-purple-500 rounded-full animate-spin" style="animation-duration: 1.5s;"></div>
                
                <!-- 中心图标：字幕/文档图标 -->
                <div class="absolute inset-0 flex items-center justify-center">
                  <svg class="w-6 h-6 text-blue-600 animate-pulse" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>
              
              <!-- 飘动的文字粒子效果 -->
              <div class="absolute inset-0 pointer-events-none">
                <div class="absolute top-2 left-8 w-1 h-1 bg-blue-400 rounded-full animate-float" style="animation-delay: 0s;"></div>
                <div class="absolute top-4 right-6 w-1 h-1 bg-indigo-400 rounded-full animate-float" style="animation-delay: 0.5s;"></div>
                <div class="absolute bottom-3 left-6 w-1 h-1 bg-purple-400 rounded-full animate-float" style="animation-delay: 1s;"></div>
                <div class="absolute bottom-2 right-8 w-1 h-1 bg-blue-400 rounded-full animate-float" style="animation-delay: 1.5s;"></div>
              </div>
            </div>
          </div>
          
          <!-- 通用旋转动画 -->
          <div v-else class="w-12 h-12 mx-auto">
            <svg class="animate-spin text-gray-600" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
          </div>
        </div>

        <!-- 进度条 -->
        <div v-if="showProgress" class="mb-6">
          <div class="w-48 bg-gray-200 rounded-full h-1.5 overflow-hidden mx-auto">
            <div 
              class="h-full bg-gradient-to-r from-gray-500 to-gray-700 rounded-full transition-all duration-1000 ease-out"
              :style="{ width: `${progress}%` }"
            ></div>
          </div>
        </div>

        <!-- 主要提示文本 -->
        <div class="text-sm font-medium text-gray-700 mb-1">
          {{ currentMessage }}
        </div>

    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed, onMounted, onUnmounted } from 'vue'

export default defineComponent({
  name: 'LoadingOverlay',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    type: {
      type: String,
      default: 'default', // 'tts', 'video', 'subtitle', 'upload', 'default'
    },
    textLength: {
      type: Number,
      default: 0
    },
    customMessages: {
      type: Array,
      default: () => []
    },
    showProgress: {
      type: Boolean,
      default: true
    }
  },
  setup(props) {
    const elapsed = ref(0)
    const progress = ref(0)
    let timer = null
    let progressTimer = null

    // 默认消息配置
    const defaultMessages = {
      tts: [
        { time: 0, message: '正在分析文本...' },
        { time: 3, message: '正在合成语音...' },
        { time: 10, message: '处理中，请稍候...' },
        { time: 30, message: '即将完成...' }
      ],
      video: [
        { time: 0, message: '正在获取视频信息...' },
        { time: 3, message: '正在解析视频内容...' },
        { time: 8, message: '即将完成...' }
      ],
      subtitle: [
        { time: 0, message: '正在提取视频字幕...' },
        { time: 8, message: '智能分析音频内容...' },
        { time: 20, message: 'AI 正在精准识别语音...' },
        { time: 40, message: '优化语言表达...' },
        { time: 60, message: '生成高质量文本...' },
        { time: 90, message: '最后的润色处理...' },
        { time: 110, message: '即将完成，请稍候...' }
      ],
      default: [
        { time: 0, message: '处理中...' },
        { time: 5, message: '请稍候...' },
        { time: 15, message: '即将完成...' }
      ]
    }

    // 根据文本长度估算时间
    const estimatedTime = computed(() => {
      if (props.type === 'tts' && props.textLength > 0) {
        // 更精确的时间预估算法
        let baseTime = 8 // 基础处理时间
        
        if (props.textLength <= 50) {
          baseTime = 5 // 短文本快速处理
        } else if (props.textLength <= 200) {
          baseTime = 8 + Math.floor(props.textLength / 50) * 2 // 每50字符增加2秒
        } else if (props.textLength <= 500) {
          baseTime = 15 + Math.floor((props.textLength - 200) / 100) * 3 // 中长文本
        } else {
          baseTime = 25 + Math.floor((props.textLength - 500) / 200) * 5 // 长文本
        }
        
        return Math.min(90, baseTime) // 最长不超过90秒
      }
      
      if (props.type === 'video') {
        return 15 // 视频信息获取通常在15秒内完成
      }
      
      if (props.type === 'subtitle') {
        return 120 // 字幕提取可能需要10-120秒
      }
      
      return 30
    })

    // 获取当前消息
    const currentMessage = computed(() => {
      const messages = props.customMessages.length > 0 
        ? props.customMessages 
        : defaultMessages[props.type] || defaultMessages.default

      // 找到当前时间对应的消息
      let current = messages[0]
      for (const msg of messages) {
        if (elapsed.value >= msg.time) {
          current = msg
        } else {
          break
        }
      }
      return current.message
    })


    // 音波高度动画
    const getBarHeight = (index) => {
      const heights = ['h-4', 'h-6', 'h-10', 'h-6', 'h-4']
      return heights[index - 1] || 'h-4'
    }

    // 开始计时
    const startTimer = () => {
      if (timer) clearInterval(timer)
      if (progressTimer) clearInterval(progressTimer)
      
      elapsed.value = 0
      progress.value = 0

      // 更新时间
      timer = setInterval(() => {
        elapsed.value += 1
      }, 1000)

      // 更新进度条（模拟进度）- 前快后慢的曲线
      progressTimer = setInterval(() => {
        if (progress.value < 90) {
          let increment = 2
          
          // 优化的速度曲线：前面快，后面慢
          if (progress.value < 20) {
            increment = 3 // 前20%很快
          } else if (progress.value < 40) {
            increment = 2 // 20-40%较快
          } else if (progress.value < 60) {
            increment = 1.2 // 40-60%中等
          } else if (progress.value < 75) {
            increment = 0.8 // 60-75%较慢
          } else if (progress.value < 85) {
            increment = 0.4 // 75-85%很慢
          } else {
            increment = 0.2 // 85%以后极慢
          }
          
          progress.value = Math.min(90, progress.value + increment)
        }
      }, 1000)
    }

    // 停止计时
    const stopTimer = () => {
      if (timer) {
        clearInterval(timer)
        timer = null
      }
      if (progressTimer) {
        clearInterval(progressTimer)
        progressTimer = null
      }
      // 完成时进度条到100%
      if (props.show === false) {
        progress.value = 100
        setTimeout(() => {
          elapsed.value = 0
          progress.value = 0
        }, 500)
      }
    }

    // 监听显示状态变化
    const handleShowChange = () => {
      if (props.show) {
        startTimer()
      } else {
        stopTimer()
      }
    }

    onMounted(() => {
      if (props.show) {
        startTimer()
      }
    })

    onUnmounted(() => {
      stopTimer()
    })

    // 监听 props.show 变化
    const unwatchShow = computed(() => props.show)
    const stopWatching = ref(null)
    
    const watchShow = () => {
      if (stopWatching.value) stopWatching.value()
      
      let oldValue = unwatchShow.value
      const checkChange = () => {
        const newValue = unwatchShow.value
        if (newValue !== oldValue) {
          handleShowChange()
          oldValue = newValue
        }
        if (!stopWatching.value) return
        requestAnimationFrame(checkChange)
      }
      
      requestAnimationFrame(checkChange)
      stopWatching.value = () => {
        stopWatching.value = null
      }
    }

    onMounted(watchShow)
    onUnmounted(() => {
      if (stopWatching.value) stopWatching.value()
    })

    return {
      elapsed,
      progress,
      estimatedTime,
      currentMessage,
      getBarHeight
    }
  }
})
</script>

<style scoped>
@keyframes wave {
  0%, 100% { 
    transform: scaleY(0.3);
    opacity: 0.7;
  }
  50% { 
    transform: scaleY(1);
    opacity: 1;
  }
}

.animate-wave {
  animation: wave 1.5s ease-in-out infinite;
  transform-origin: bottom;
}

/* 飘动的粒子动画 */
@keyframes float {
  0%, 100% { 
    transform: translateY(0px) translateX(0px);
    opacity: 0.3;
  }
  25% { 
    transform: translateY(-8px) translateX(4px);
    opacity: 1;
  }
  50% { 
    transform: translateY(-12px) translateX(-2px);
    opacity: 0.8;
  }
  75% { 
    transform: translateY(-6px) translateX(3px);
    opacity: 0.6;
  }
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

</style>
