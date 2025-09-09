<template>
  <div class="flex flex-col gap-2">
    <label v-if="label" class="text-sm text-gray-600">{{ label }}</label>
    <textarea
      :readonly="readonly"
      v-model="localValue"
      class="w-full h-64 p-3 border rounded-md  focus:outline-none focus:ring"
      :placeholder="placeholder"
      @input="onInput"
    />
  </div>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'TextEditor',
  props: {
    modelValue: { type: String, default: '' },
    placeholder: { type: String, default: '' },
    label: { type: String, default: '' },
    readonly: { type: Boolean, default: false },
  },
  emits: ['update:modelValue'],
  data() {
    return { localValue: this.modelValue }
  },
  watch: {
    modelValue(v) {
      if (v !== this.localValue) this.localValue = v
    },
  },
  methods: {
    onInput() {
      this.$emit('update:modelValue', this.localValue)
    },
  },
})
</script>




