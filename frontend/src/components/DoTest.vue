<template>
  <div class="container">
    <div v-if="!allowed">
      <br>
      <text-input v-model="client_test.jmbg" :required=true label="Upišite vaš JMBG" name="jmbg" type="text">
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

          <h6>Pitanja</h6>
          <div v-for="(question, index) in questions" :key="question.ID" class="row mt-2" style="margin-bottom: 3px">
            <div class="col-xs-12 col-sm-12 col-md-12">
              <p> {{ index + 1 }}. {{ question.content }}</p>
              <img v-if="question.image" id="img" :src="question.image" alt=""
                style="margin-bottom: 3px; max-width: 100%; height: auto" />
              <div v-for="(answer) in question.answers" :key="answer.ID" class="row no-gutters">
                <div class="col-2 col-sm-2 col-md-1"
                  :style="[(isSecondFinished && answer.correct) ? { 'background': 'green' } : {}]">
                  {{ answer.letter }}
                  <input style="margin-left: 5px" id={{question.ID}} v-model="client_test.questions_answers[index].answer"
                    :value=answer.letter type="radio">
                </div>
                <div class="col-10 col-sm-10 col-md-10">{{ answer.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <input v-if="!isfinish" :disabled="disableSaveButton" class="btn btn-primary m-2" type="submit" value="Snimi">

      <div>
        <div v-if="isfinish">
          <p v-if="readyForSurvey" class="bg-secondary mt-1">Molimo vas da popunite anketu u nastavku.</p>
          <p class="bg-info mt-1">Rezultat: {{ (client_test.result * 100).toFixed(2) }}%</p>
        </div>

        <div v-if="readyForSurvey && surveyDataReady">
          <h6>Anketa</h6>
          <p>1 - у потпуности се не слажем, 5 - у потпуности се слажем </p> 
          <div v-for="(question, index) in surveys[0].questions" :key="question.ID" class="row mt-2"
            style="margin-bottom: 3px">
            <div class="col-xs-12 col-sm-12 col-md-12">
              <p> {{ index + 1 }}. {{ question.content }}</p>
              <div class="col-6 col-sm-12 col-md-6">
                1<input style="margin-left: 5px" id={{question.ID}}
                  v-model="clientSurveys[0].survey_questions_answers[index].grade" :value=1 type="radio">
                2<input style="margin-left: 5px" id={{question.ID}}
                  v-model="clientSurveys[0].survey_questions_answers[index].grade" :value=2 type="radio">
                3<input style="margin-left: 5px" id={{question.ID}}
                  v-model="clientSurveys[0].survey_questions_answers[index].grade" :value=3 type="radio">
                4<input style="margin-left: 5px" id={{question.ID}}
                  v-model="clientSurveys[0].survey_questions_answers[index].grade" :value=4 type="radio">
                5<input style="margin-left: 5px" id={{question.ID}}
                  v-model="clientSurveys[0].survey_questions_answers[index].grade" :value=5 type="radio">
              </div>
            </div>
          </div>

          <hr>
          <h6>Anketa o predavačima</h6>

          <div v-for="(teacher, i) in teachers" :key="teacher.ID" class="row mt-2" style="margin-bottom: 3px">
            <hr>
            <p style="font-weight:bold;">{{ teacher.person.first_name }} {{ teacher.person.last_name }}</p>
            <div v-for="(question, index) in surveys[1].questions" :key="question.ID" class="row mt-2"
              style="margin-bottom: 3px">
              <div class="col-xs-12 col-sm-12 col-md-12">
                <p> {{ index + 1 }}. {{ question.content }}</p>
                <div class="col-6 col-sm-12 col-md-6">
                  1<input style="margin-left: 5px" id={{question.ID}}
                    v-model="clientSurveys[i + 1].survey_questions_answers[index].grade" :value=1 type="radio">
                  2<input style="margin-left: 5px" id={{question.ID}}
                    v-model="clientSurveys[i + 1].survey_questions_answers[index].grade" :value=2 type="radio">
                  3<input style="margin-left: 5px" id={{question.ID}}
                    v-model="clientSurveys[i + 1].survey_questions_answers[index].grade" :value=3 type="radio">
                  4<input style="margin-left: 5px" id={{question.ID}}
                    v-model="clientSurveys[i + 1].survey_questions_answers[index].grade" :value=4 type="radio">
                  5<input style="margin-left: 5px" id={{question.ID}}
                    v-model="clientSurveys[i + 1].survey_questions_answers[index].grade" :value=5 type="radio">
                </div>
              </div>
            </div>
          </div>

          <input :disabled="disableSaveSurveyButton" class="btn btn-primary m-2" value="Snimi"
            @click.prevent="createClientSurvey()">
        </div>

      </div>
    </form-tag>
  </div>
</template>

<script>
import FormTag from "@/components/forms/FormTag";
import axios from "axios";
import { apiMixin } from "@/mixins/apiMixin";
import { styleMixin } from "@/mixins/styleMixin";
import { useToast } from "vue-toastification";
import TextInput from "@/components/forms/TextInput.vue";
import { dateMixin } from "@/mixins/dateMixin";
import { commonMixin } from "@/mixins/commonMixin";
import { fileMixin } from "@/mixins/fileMixin";

export default {
  name: 'DoTest',
  mixins: [apiMixin, styleMixin, dateMixin, commonMixin, fileMixin],
  components: { TextInput, FormTag },
  data() {
    return {
      client_test: {
        jmbg: "",
        seminar_day: null,
        seminar_day_id: 0,
        test: null,
        questions_answers: [],
        result: 0.0,
        isSecondFinished: false,
      },
      seminarDay: {
        ID: 0,
        date: null,
        test_id: 0,
        test: null,
      },
      questions: [],
      surveys: [],
      teachers: [],
      seminarDayId: 0,
      isfinish: false,
      surveyDataReady: false,
      isSecondFinished: false,
      allowed: false,
      disableSaveButton: false,
      disableSaveSurveyButton: false,
      disableJMBGButton: false,
      clientSurveys: [],
    }
  },
  computed: {
    readyForSurvey() {
      return this.isfinish && this.isSecondFinished;
    },
  },
  methods: {
    async createClientSurvey() {
      this.disableSaveSurveyButton = true;
      var showSuccess = true;
      this.clientSurveys.forEach(cs => {
      cs.seminar_day = this.client_test.seminar_day;
      cs.seminar_day_id = this.client_test.seminar_day_id;
      axios.post('/surveys/client-survey/create', JSON.stringify(cs)).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          showSuccess = false;
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
      }, (error) => {
        showSuccess = false;
        this.errorToast(error, "/surveys/client-survey/create");
      });
      });
      if (showSuccess) {
        this.toast.info("Uspešno snimljena anketa!");
      }
    },
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
        this.seminarDay.date = this.getDateInMMDDYYYYFormat(this.seminarDay.date);

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
        this.questions.sort((a, b) => (a - b));
        this.questions.forEach(q => {
          q.answers = q.answers.sort((a, b) => {
            if (a.letter < b.letter) {
              return -1;
            }
            if (a.letter > b.letter) {
              return 1;
            }
            return 0;
          });
          const obj = { question_id: q.ID, answer: "" };
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

      if (this.readyForSurvey) {
        await this.getTeachersBySeminarDay();
        await this.getActiveSurveys();
      }
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
        this.isSecondFinished = newClientTest.isSecondFinished;
        this.toast.info("Uspešno snimljen test!");
      }, (error) => {
        this.errorToast(error, "/tests/client-test/create");
      });
    },
    async getActiveSurveys() {
      await axios.get('/surveys/active').then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        const result = JSON.parse(response.data.Data);

        var generalSurvey = result.find((s) => s.type == this.SURVEY_TYPES.GENERAL);
        this.surveys.push(generalSurvey);
        var teacherSurvey = result.find((s) => s.type == this.SURVEY_TYPES.TEACHER);
        this.surveys.push(teacherSurvey);

        //general survey
        var generalClientSurvey = {}
        generalClientSurvey.jmbg = this.client_test.jmbg;
        generalClientSurvey.survey = generalSurvey;
        generalClientSurvey.survey_id = generalSurvey.id;
        var question_answers = [];
        generalSurvey.questions.forEach(q => {
          const obj = { survey_question_id: q.ID, grade: 0 };
          question_answers.push(obj);
        });

        generalClientSurvey.survey_questions_answers = question_answers;

        this.clientSurveys.push(generalClientSurvey);

        //teacher surveys
        this.teachers.forEach(t => {
          var clientSurvey = {};
          clientSurvey.jmbg = this.client_test.jmbg;
          clientSurvey.survey = teacherSurvey;
          clientSurvey.survey_id = teacherSurvey.id;
          clientSurvey.teacher_id = t.id;
          clientSurvey.teacher = t;
          var question_answers = [];
          teacherSurvey.questions.forEach(q => {
            const obj = { survey_question_id: q.ID, grade: 0 };
            question_answers.push(obj);

          });
          clientSurvey.survey_questions_answers = question_answers;

          this.clientSurveys.push(clientSurvey);
        });
        console.log(this.clientSurveys);
        this.surveyDataReady = true;
      }, (error) => {
        this.errorToast(error, "/surveys/active");
      });
    },
    async getTeachersBySeminarDay() {
      await axios.get('/seminar-days/teachers/id/' + this.client_test.seminar_day.ID.toString()).then((response) => {
        if (response.data === null || response.data.Status === 'error') {
          this.toast.error(response.data != null ? response.data.ErrorMessage : "");
          return;
        }
        this.teachers = JSON.parse(response.data.Data);
        console.log(this.teachers);
      }, (error) => {
        this.errorToast(error, "/seminar-days/teachers/id/");
      });
    }
  },
  setup() {
    const toast = useToast();
    return { toast }
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
}</style>