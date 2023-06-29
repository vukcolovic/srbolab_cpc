<template>
  <div class="mb-1">
    <label :for="name" class="form-label" :style="styleLabel">{{label}}</label>
    <textarea
        :name="name"
        :placeholder="placeholder"
        :readonly="readonly"
        :required="required"
        :min="min"
        :max="max"
        :rows="rows"
        :value="modelValue"
        @input="$emit('update:modelValue', $event.target.value)"
        class="form-control"
        :style="styleInput">
    </textarea>
  </div>
</template>

<script>
export default {
  name: 'TextArea',
  props: {
    name: String,
    label: String,
    placeholder: String,
    required: String,
    min: String,
    max: String,
    rows: String,
    modelValue: String,
    readonly: Boolean,
    styleInput: String,
    styleLabel: String,
  },
}
</script>

<style scoped>


</style>