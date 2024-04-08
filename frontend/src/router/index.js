import { createRouter, createWebHistory } from 'vue-router'
import BodyComponent from './../components/BodyComponent.vue'
import LoginComponent from './../components/LoginComponent.vue'
import ClientWaitingList from '../components/ClientWaitingList.vue'
import ClientsList from '../components/ClientsList.vue'
import ClientEdit from './../components/ClientEdit.vue'
import UsersList from '../components/UsersList.vue'
import UserEdit from './../components/UserEdit.vue'
import SupportComponent from './../components/SupportComponent.vue'
import SeminarsList from '../components/SeminarsList.vue'
import SeminarEdit from './../components/SeminarEdit.vue'
import CompaniesList from '../components/CompaniesList.vue'
import CompanyEdit from './../components/CompanyEdit.vue'
import LocationsList from '../components/LocationsList.vue'
import LocationEdit from './../components/LocationEdit.vue'
import ClassRoomsList from '../components/ClassRoomsList.vue'
import ClassRoomEdit from './../components/ClassRoomEdit.vue'
import SeminarDayEdit from './../components/SeminarDayEdit.vue'
import QuestionsList from '../components/QuestionsList.vue'
import QuestionEdit from '../components/QuestionEdit.vue'
import TestsList from '../components/TestsList.vue'
import TestEdit from '../components/TestEdit.vue'
import ClientEditNoCorporate from '../components/ClientEditNoCorporate.vue'
import DoTest from '../components/DoTest.vue'
import ClassNamesList from '../components/ClassNamesList.vue'
import ClassNameEdit from '../components/ClassNameEdit.vue'
import ThemeNamesList from '../components/ThemeNamesList.vue'
import ThemeNameEdit from '../components/ThemeNameEdit.vue'
import SurveyQuestionsList from '../components/SurveyQuestionsList.vue'
import SurveyQuestionEdit from '../components/SurveyQuestionEdit.vue'
import SurveysList from '../components/SurveysList.vue'
import SurveyEdit from '../components/SurveyEdit.vue'
import PartnersList from '../components/PartnersList.vue'
import PartnerEdit from '../components/PartnerEdit.vue'

const routes = [
  {
    path: '/',
    name: 'ClientEditNoCorporate',
    component: ClientEditNoCorporate,
  },
  {
    path: '/do-test',
    name: 'DoTest',
    component: DoTest,
    props: true,
  },
  {
    path: '/home',
    name: 'Home',
    component: BodyComponent,
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginComponent,
  },
  {
    path: '/support',
    name: 'Support',
    component: SupportComponent,
  },
  {
    path: '/clients-waiting',
    name: 'ClientWaitingList',
    component: ClientWaitingList,
  },
  {
    path: '/clients',
    name: 'ClientsList',
    component: ClientsList,
  },
  {
    path: '/client',
    name: 'ClientEdit',
    component: ClientEdit,
    props: true,
  },
  {
    path: '/seminars',
    name: 'Seminars',
    component: SeminarsList,
  },
  {
    path: '/seminar',
    name: 'SeminarEdit',
    component: SeminarEdit,
    props: true,
  },
  {
    path: '/users',
    name: 'Users',
    component: UsersList,
  },
  {
    path: '/user',
    name: 'UserEdit',
    component: UserEdit,
    props: true,
  },
  {
    path: '/companies',
    name: 'Companies',
    component: CompaniesList,
  },
  {
    path: '/company',
    name: 'CompanyEdit',
    component: CompanyEdit,
    props: true,
  },
  {
    path: '/locations',
    name: 'Locations',
    component: LocationsList,
  },
  {
    path: '/location',
    name: 'LocationEdit',
    component: LocationEdit,
    props: true,
  },
  {
    path: '/class-rooms',
    name: 'classRooms',
    component: ClassRoomsList,
  },
  {
    path: '/class-room',
    name: 'ClassRoomEdit',
    component: ClassRoomEdit,
    props: true,
  },
  {
    path: '/seminar-day',
    name: 'SeminarDayEdit',
    component: SeminarDayEdit,
    props: true,
  },
  {
    path: '/questions',
    name: 'Questions',
    component: QuestionsList,
  },
  {
    path: '/question',
    name: 'QuestionEdit',
    component: QuestionEdit,
    props: true,
  },
  {
    path: '/tests',
    name: 'Tests',
    component: TestsList,
  },
  {
    path: '/test',
    name: 'TestEdit',
    component: TestEdit,
    props: true,
  },
  {
    path: '/class-names',
    name: 'ClassNames',
    component: ClassNamesList,
  },
  {
    path: '/class-name',
    name: 'ClassNameEdit',
    component: ClassNameEdit,
    props: true,
  },
  {
    path: '/theme-names',
    name: 'ThemeNames',
    component: ThemeNamesList,
  },
  {
    path: '/theme-name',
    name: 'ThemeNameEdit',
    component: ThemeNameEdit,
    props: true,
  },
  {
    path: '/survey-questions',
    name: 'SurveyQuestions',
    component: SurveyQuestionsList,
  },
  {
    path: '/survey-question',
    name: 'SurveyQuestionEdit',
    component: SurveyQuestionEdit,
    props: true,
  },
  {
    path: '/surveys',
    name: 'Surveys',
    component: SurveysList,
  },
  {
    path: '/survey',
    name: 'SurveyEdit',
    component: SurveyEdit,
    props: true,
  },
  {
    path: '/partners',
    name: 'Partners',
    component: PartnersList,
  },
  {
    path: '/partner',
    name: 'PartnerEdit',
    component: PartnerEdit,
    props: true,
  },
]

const router = createRouter({history: createWebHistory(), routes})
export default router