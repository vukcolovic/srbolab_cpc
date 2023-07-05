import axios from "axios";
import {useToast} from "vue-toastification";

export const apiMixin = {
    data() {
        return {
            locations: [],
        }
    },
    methods: {
        async getAllLocations() {
            await axios.get('/locations/list').then((response) => {
                if (response.data === null || response.data.Status === 'error') {
                    this.toast.error(response.data.ErrorMessage);
                }
                this.locations = JSON.parse(response.data.Data);
            }, (error) => {
                this.toast.error(error);
            });
        },
    },
    setup() {
        const toast = useToast();
        return {toast}
    },
}