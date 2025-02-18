export const useForm = (initFormData = {}) => {
  const formRef = ref(null)
  const formModel = ref({ ...initFormData })
  const rules = {
    required: {
      required: true,
      message: '此为必填项',
      trigger: ['blur', 'change'],
    },
  }
  const validation = () => {
    return formRef.value?.validate()
  }
  return [formRef, formModel, validation, rules]
}
