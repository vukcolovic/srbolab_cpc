<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-11 mx-auto">
        <button class="btn btn-primary" @click="backToSeminar()"><i class="fa fa-arrow-left"></i>Vrati se na seminar</button>
      </div>
    </div>
    <form-tag @formEvent="submitHandler" name="myForm" event="formEvent">
      <div class="row">
        <div class="col-sm-5">
          <label style="font-size: 2em">{{seminarDay.seminar_theme}} - Dan {{seminarDay.number}}</label>
          <text-input
              v-model.trim="seminarDay.name"
              label="Naziv"
              type="text"
              name="name"
              :required=true
              :readonly="readonly">
          </text-input>

          <text-input
              v-model="seminarDay.date"
              label="Datum"
              type="date"
              name="date"
              :required=true
              :readonly="readonly">
          </text-input>


          <input type="submit" class="btn btn-primary m-2" value="Snimi">
        </div>

      </div>
    </form-tag>
  </div>
</template>

<script>
import TextInput from "@/components/forms/TextInput";
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import router from "@/router";
import {fileMixin} from "@/mixins/fileMixin";
import {useToast} from "vue-toastification";
import {apiMixin} from "@/mixins/apiMixin";

export default {
  name: 'SeminarDayEdit',
  mixins: [fileMixin, apiMixin],
  components: {FormTag, TextInput},
  computed: {
    readonly() {
      return this.action === 'view';
    },
  },
  data() {
    return {
      seminarDay: {ID: 0, date: null, number: 0, name: "", seminar_id: 0, seminar: null, seminar_theme: ""},
      seminarDayId: 0,
    }
  },
  methods: {
    backToSeminar(){
      router.push({name: 'SeminarEdit', query: {id: this.seminarDay.seminar_id, action: "update"}});
    },
    async getSeminarDayById() {
      axios.get('/seminar-days/id/' + this.seminarDayId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarDay = JSON.parse(response.data.Data);
        this.seminarDay.date = this.getDateInMMDDYYYYFormat(this.seminarDay.date);
        this.seminarDay.seminar_theme = this.getSeminarFullType(this.seminarDay.seminar.seminar_theme.base_seminar_type, this.seminarDay.seminar.seminar_theme.name);
      }, (error) => {
        this.toast.error(error.message);
      });
    },
    async submitHandler() {
      await axios.post('/seminar-days/update', JSON.stringify(this.seminarDay)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.toast.info("Uspešno ažuriran seminar dan.");
      }, (error) => {
        this.toast.error(error.message);
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    this.seminarDayId = this.$route.query.id;
    await this.getSeminarDayById();
  }
}
</script>