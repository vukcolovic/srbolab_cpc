import {useToast} from "vue-toastification";

export const commonMixin = {
    data() {
        return {
            yesNoOptions: [{label: 'ДА', value: 'true'}, {label: 'НЕ', value: 'false'}, {label: '-', value: ''}],
        }
    },
    setup() {
        const toast = useToast();
        return {toast}
    },
}