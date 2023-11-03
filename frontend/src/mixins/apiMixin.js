import axios from "axios";
import {useToast} from "vue-toastification";
import {dateMixin} from "@/mixins/dateMixin";
import router from "@/router";

export const apiMixin = {
    mixins: [dateMixin],
    data() {
        return {
            locations: [],
            seminarBaseTypes: [],
            seminarThemes: [],
            seminarStatuses: [],
            companies: [],
            users: [],
        }
    },
    methods: {
        async getAllUsers() {
            await axios.get('/users/list?skip=' + 0 + '&take=' + 1000).then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                    return;
                }
                this.users = JSON.parse(response.data.Data);
                this.users.forEach(vs => {
                    vs.first_name = vs.person.first_name;
                    vs.last_name = vs.person.last_name;
                    vs.email = vs.person.email;
                    vs.phone_number = vs.person.phone_number;
                    vs.full_name = vs.person.first_name + " " + vs.person.last_name;
                });
            }, (error) => {
                this.errorToast(error, "/users/list");
            });
        },
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
                this.errorToast(error, "/locations/list");
            });
        },
        async getAllCompanies() {
            await axios.get('/companies/list?skip=0&take=10000').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                    return;
                }
                this.companies = JSON.parse(response.data.Data);
                this.companies.forEach(com => {
                    com.name_pib = com.name + "-" + com.pib;
                });
            }, (error) => {
                this.errorToast(error, "/companies/list");
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
                this.errorToast(error, "/seminar-types/list");
            });
        },
        async getAllSeminarThemes() {
            await axios.get('/seminar-types/themes/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                    return;
                }
                this.seminarThemes = JSON.parse(response.data.Data);
            }, (error) => {
                this.errorToast(error, "/seminar-types/themes/list");
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
                this.errorToast(error, "/seminar-statuses/list");
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
                    vs.base_info = this.getSeminarFullType(vs.seminar_theme.base_seminar_type, vs.seminar_theme) + " | " + this.getDateInMMDDYYYYFormat(vs.start_date) + " | " + vs.class_room.location.code + "-" + vs.class_room.name;
                });
                return result;
            }, (error) => {
                this.errorToast(error, "/seminars/list/status/");
            });
        },
        async isCorporateIp() {
            await axios.get('/corporate-ip').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    return;
                }
                var corporate = JSON.parse(response.data.Data);
                this.$store.dispatch('setIsCorporateAction', corporate.is_corporate);
                if (corporate.is_corporate) {
                    router.push("/home");
                }
            }, (error) => {
                this.errorToast(error, "/corporate-ip");
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