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
            companies: [],
        }
    },
    methods: {
        async getAllLocations() {
            await axios.get('/locations/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                    return;
                }
                this.locations = JSON.parse(response.data.Data);
                this.locations.forEach(s => {
                    s.address_place = s.address.place;
                });
            }, (error) => {
                this.toast.error(error.message);
            });
        },
        async getAllCompanies() {
            await axios.get('/companies/list?skip=0&take=10000').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                    return;
                }
                this.companies = JSON.parse(response.data.Data);
            }, (error) => {
                this.toast.error(error.message);
            });
        },
        async getAllBaseSeminarTypes() {
            await axios.get('/seminar-types/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                    return;
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
                    return;
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
                    vs.base_info = this.getSeminarFullType(vs.seminar_theme.base_seminar_type, vs.seminar_theme) + " | " + this.getDateInMMDDYYYYFormat(vs.start_date) + " | " + vs.class_room.location.address.place;
                });
                return result;
            }, (error) => {
                this.toast.error(error.message);
            });
        },
        getSeminarFullType(base, theme) {
            if (base.code === "ADDITIONAL" || base.code === "BASIC") {
                return base.name;
            }

            return base.name + " " + theme.name;
        }
    },
    setup() {
        const toast = useToast();
        return {toast}
    },
}