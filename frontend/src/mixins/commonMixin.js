import {useToast} from "vue-toastification";

const SEMINAR_STATUSES = {
    OPENED: 1,
    FILLED: 2,
    IN_PROGRESS: 3,
    CLOSED: 4
};

export const commonMixin = {
    data() {
        return {
            SEMINAR_STATUSES: SEMINAR_STATUSES,
            yesNoOptions: [{label: 'ДА', value: 'true'}, {label: 'НЕ', value: 'false'}, {label: '-', value: ''}],
        }
    },
    setup() {
        const toast = useToast();
        return {toast}
    },
}