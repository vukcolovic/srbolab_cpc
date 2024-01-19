import {useToast} from "vue-toastification";

const SEMINAR_STATUSES = {
    OPENED: 1,
    FILLED: 2,
    IN_PROGRESS: 3,
    CLOSED: 4
};

const SURVEY_TYPES = {
    GENERAL: 0,
    TEACHER: 1,
};

export const commonMixin = {
    data() {
        return {
            SEMINAR_STATUSES: SEMINAR_STATUSES,
            SURVEY_TYPES: SURVEY_TYPES,
            yesNoOptions: [{label: 'ДА', value: 'true'}, {label: 'НЕ', value: 'false'}],
        }
    },
    methods: {
        async errorToast(error, api) {
            if (error == null) {
                this.toast.error("Greška prilikom poziva " + api);
                return;
            }
            if (error.data) {
                this.toast.error(error.data);
                return;
            }
            this.toast.error(error);
        },
    },
    setup() {
        const toast = useToast();
        return {toast}
    },
}