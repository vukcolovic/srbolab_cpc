import axios from "axios";
import {useToast} from "vue-toastification";
import {dateMixin} from "@/mixins/dateMixin";

export const apiMixin = {
    mixins: [dateMixin],
    data() {
        return {
            locations: [],
            seminarBaseTypes: [],
            seminarStatuses: [],
        }
    },
    methods: {
        async getAllLocations() {
            await axios.get('/locations/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                }
                this.locations = JSON.parse(response.data.Data);
                this.locations.forEach(s => {
                    s.address_place = s.address.place;
                });
            }, (error) => {
                this.toast.error(error.message);
            });
        },
        async getAllBaseSeminarTypes() {
            await axios.get('/seminar-types/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                }
                this.seminarBaseTypes = JSON.parse(response.data.Data);
            }, (error) => {
                this.toast.error(error.message);
            });
        },
        async getAllSeminarStatuses() {
            await axios.get('/seminar-statuses/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                }
                this.seminarStatuses = JSON.parse(response.data.Data);
            }, (error) => {
                this.toast.error(error.message);
            });
        },
        async getSeminarsByStatusCode(statusCode) {
            return await axios.get('/seminars/list/status/' + statusCode).then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                    return;
                }
                var result = JSON.parse(response.data.Data);
                result.forEach(vs => {
                    vs.details = vs.seminar_theme.base_seminar_type.name + "-" + vs.seminar_theme.name + " | " + this.getDateInMMDDYYYYFormat(vs.start_date) + " | " + vs.class_room.location.address.place;
                });
                return result;
            }, (error) => {
                this.toast.error(error.message);
            });
        },
    },
    setup() {
        const toast = useToast();
        return {toast}
    },
}