<template>
  <div class="container">
    <form-tag event="formEvent" name="myForm" @formEvent="submitHandler">
      <div class="row">
        <div class="col-sm-11 mx-auto">
          <h3 class="mt-2">Test - {{ seminarDay.test ? seminarDay.test.name : "" }}</h3>
          <hr>
        </div>
        <div class="col-sm-12">
          <text-input
              v-model="client_test.jmbg"
              :required=true
              label="Upišite vaš JMBG"
              name="jmbg"
              type="text">
          </text-input>
          <hr>
          <h6>Pitanja</h6>
          <div v-for="(question, index) in questions" :key="question.ID" class="row" style="margin-bottom: 5px">
            <div class="col-sm-11">
              <p> {{ index + 1 }}. {{ question.content }}</p>
              <div v-for="(answer) in question.answers" :key="answer.ID" class="row">
                {{ answer.letter }}) <input id={{question.ID}} v-model="client_test.questions_answers[index].answer" :value= answer.letter
                                            class="col-sm-1" type="radio">
                {{ answer.content }}
                <p></p>
              </div>
            </div>
          </div>

        </div>
      </div>
      <div>
        <input class="btn btn-primary m-2" type="submit" value="Snimi">
        <p v-if="isfinish" class="bg-info mt-1">Rezultat: {{client_test.result * 100}}%</p>
      </div>
    </form-tag>
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
        test: null,
        questions_answers: [],
        result: 0.0,
      },
      seminarDay: {
        ID: 0,
        test_id: 0,
        test: null,
      },
      questions: [],
      seminarDayId: 0,
      isfinish: false,
    }
  },
  methods: {
    async getSeminarDayById() {
      await axios.get('/seminar-days/id/' + this.seminarDayId).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.seminarDay = JSON.parse(response.data.Data);
        this.questions = this.seminarDay.test.questions;
        this.questions.sort((a, b) => (a - b))
        this.questions.forEach(q => {
          const obj = {question_id: q.ID, answer: ""};
          this.client_test.questions_answers.push(obj);
        })

      }, (error) => {
        this.toast.error(error.message);
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
      this.client_test.seminar_day = this.seminarDay;
      this.client_test.test = this.seminarDay.test;
      await axios.post('/tests/client-test/create', JSON.stringify(this.client_test)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        const newClientTest = JSON.parse(response.data.Data);
        this.client_test.result = newClientTest.result;
        this.isfinish = true;
        this.toast.info("Uspešno snimljen test!");
      }, (error) => {
        this.toast.error(error ? error.data : "");
      });
    },
  },
  setup() {
    const toast = useToast();
    return {toast}
  },
  async mounted() {
    if (this.$route.query.seminar_day_id === '') {
      this.toast.error("URL nije validan!");
      return;
    }
    this.seminarDayId = this.$route.query.seminar_day_id;
    await this.getSeminarDayById();
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