import axios from "axios";
import {useToast} from "vue-toastification";

export const apiMixin = {
    data() {
        return {
            locations: [],
            seminarTypes: [],
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
        async getAllSeminarTypes() {
            await axios.get('/seminar-types/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data != null ? response.data.ErrorMessage : "");
                }
                this.seminarTypes = JSON.parse(response.data.Data);
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
    },
    setup() {
        const toast = useToast();
        return {toast}
    },
}