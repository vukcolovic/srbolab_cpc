<template>
    <div :class="[horizontal ? 'row' : 'form-label', 'mb-0', 'mt-0']">
        <label :for="name" :class="[horizontal ? horizontalLabelClasses : 'form-label']" :style="styleLabel">{{label}}</label>
      <div :class="[horizontal ? horizontalInputClasses : '']">
        <input
            :style="styleInput"
            :type="type"
            :name="name"
            :placeholder="placeholder"
            :readonly="readonly"
            :required="required"
            :step= "step"
            :min="min"
            :max="max"
            :value="modelValue"
            :autocomplete="name + '-new'"
            @input="$emit('update:modelValue', $event.target.value)"
            :ref="refVal"
            class="form-control">
      </div>
    </div>
</template>

<script>
export default {
    name: 'TextInput',
    props: {
        name: String,
        type: String,
        label: String,
        placeholder: String,
        required: Boolean,
        min: String,
        max: String,
        step: String,
        modelValue: String,
        readonly: Boolean,
        refVal: String,
        horizontal: Boolean,
        styleInput: String,
        styleLabel: String,
        horizontalLabelClasses: {
          default: 'col-sm-3 col-form-label',
          type: String
        },
      horizontalInputClasses: {
        default: 'col-sm-6',
        type: String
      }
    },
}
</script>
