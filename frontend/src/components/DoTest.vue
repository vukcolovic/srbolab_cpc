<template>
  <div class="container">
    <div v-if="!allowed">
      <br>
      <text-input
          v-model="client_test.jmbg"
          :required=true
          label="Upišite vaš JMBG"
          name="jmbg"
          type="text">
      </text-input>
      <br>
      <button type="button" :disabled="disableJMBGButton" @click="sendJmbg()" class="btn btn-primary">Pošalji</button>
    </div>

    <form-tag v-if="allowed" event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-11 mx-auto">
          <h3 class="mt-2">Test - {{ seminarDay.test ? seminarDay.test.name : "" }}</h3>
          <hr>
        </div>
        <div class="col-sm-12">
<!--          <text-input-->
<!--              v-model="client_test.jmbg"-->
<!--              :required=true-->
<!--              label="Upišite vaš JMBG"-->
<!--              name="jmbg"-->
<!--              type="text">-->
<!--          </text-input>-->
          <h6>Pitanja</h6>
          <div v-for="(question, index) in questions" :key="question.ID" class="row mt-2" style="margin-bottom: 3px">
            <div class="col-xs-12 col-sm-12 col-md-12">
              <p> {{ index + 1 }}. {{ question.content }}</p>
              <img v-if="question.image" id="img"  :src="question.image" alt="" style="margin-bottom: 3px; max-width: 100%; height: auto" />
              <div v-for="(answer) in question.answers" :key="answer.ID" class="row no-gutters">
                <div class="col-2 col-sm-2 col-md-1">
                  {{ answer.letter }})
                  <input style="margin-left: 5px" id={{question.ID}} v-model="client_test.questions_answers[index].answer" :value= answer.letter type="radio">
                </div>
                <div class="col-10 col-sm-10 col-md-10">{{ answer.content }}</div>
              </div>
            </div>
          </div>

        </div>
      </div>
      <div>
        <p v-if="isfinish" class="bg-info mt-1">Rezultat: {{(client_test.result * 100).toFixed(2)}}%</p>
        <input v-if="!isfinish" :disabled="disableSaveButton" class="btn btn-primary m-2" type="submit" value="Snimi">
      </div>
    </form-tag>
<!--    <h2 v-else>Test nije dozvoljen, obratite se rukovodiocu kursa.</h2>-->
  </div>
</template>

<script>
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import {apiMixin} from "@/mixins/apiMixin";
import {styleMixin} from "@/mixins/styleMixin";
import {useToast} from "vue-toastification";
import TextInput from "@/components/forms/TextInput.vue";
import {dateMixin} from "@/mixins/dateMixin";
import {commonMixin} from "@/mixins/commonMixin";
import {fileMixin} from "@/mixins/fileMixin";

export default {
  name: 'DoTest',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: {TextInput, FormTag},
  data() {
    return {
      client_test: {
        jmbg: "",
        seminar_day: null,
        seminar_day_id: 0,
        test: null,
        questions_answers: [],
        result: 0.0,
      },
      seminarDay: {
        ID: 0,
        date: null,
        test_id: 0,
        test: null,
      },
      questions: [],
      seminarDayId: 0,
      isfinish: false,
      allowed: false,
      disableSaveButton: false,
      disableJMBGButton: false,
    }
  },
  methods: {
    async sendJmbg() {
      if (this.disableJMBGButton) {
        return;
      }
      this.disableJMBGButton = true;
      await axios.get('/seminar-days/jmbg/' + this.client_test.jmbg).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }

        this.seminarDay = JSON.parse(response.data.Data);
        this.seminarDay.date= this.getDateInMMDDYYYYFormat(this.seminarDay.date);

        if (!this.isToday(new Date(this.seminarDay.date))) {
          this.toast.warning("Ovaj test danas nije dozvoljen!");
          return;
        }
        if (this.seminarDay.seminar.seminar_status_id != this.SEMINAR_STATUSES.IN_PROGRESS) {
          this.toast.warning("Seminar nije u toku, test ne može da se radi.");
          return;
        }
        if (!this.seminarDay.test_id) {
          this.toast.warning("Za ovaj seminar nije odabran test.");
          return;
        }
        this.allowed = true;
        this.questions = this.seminarDay.test.questions;
        this.questions.sort((a, b) => (a - b))
        this.questions.forEach(q => {
          const obj = {question_id: q.ID, answer: ""};
          this.client_test.questions_answers.push(obj);
        })

      }, (error) => {
        this.errorToast(error, "/seminar-days/jmbg");
      });
    },
    async submitHandler() {
      if (this.client_test.jmbg.length !== 13) {
        this.toast.warning("JMBG nije validan, molim vas proverite unos.");
        return;
      }
      await this.createClientTest();
    },
    async createClientTest() {
      this.disableSaveButton = true;
      this.seminarDay.date = this.getBackendFormat(this.seminarDay.date);
      delete this.seminarDay.test.questions.forEach(q => q.image = "")
      this.client_test.seminar_day = this.seminarDay;
      this.client_test.test = this.seminarDay.test;
      await axios.post('/tests/client-test/create', JSON.stringify(this.client_test)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.seminarDay.date = this.getDateInMMDDYYYYFormat(this.seminarDay.date);
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        const newClientTest = JSON.parse(response.data.Data);
        this.client_test.result = newClientTest.result;
        this.isfinish = true;
        this.toast.info("Uspešno snimljen test!");
      }, (error) => {
        this.errorToast(error, "/tests/client-test/create");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {

  }
}
</script>

<style scoped>
.shell {
  height: 20px;
  width: 500px;
  border: 1px solid #aaa;
  border-radius: 13px;
}

.bar {
  background: linear-gradient(to right, #007bff, #007bff);
  height: 20px;
  width: 15px;
  border-radius: 9px;

  span {
    float: right;
    color: #fff;
    font-size: 0.7em
  }
}

.styled-table {
  border-collapse: collapse;
  margin: 0px 0;
  font-family: sans-serif;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
}

.styled-table thead tr {
  background-color: #009879;
  color: #ffffff;
  text-align: left;
}

.styled-table th,
.styled-table td {
  padding: 1px 1px;
}

.styled-table tbody tr {
  border-bottom: 1px solid #dddddd;
}

.styled-table tbody tr:nth-of-type(even) {
  background-color: #f3f3f3;
}
</style>