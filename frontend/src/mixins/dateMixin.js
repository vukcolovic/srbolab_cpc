export const dateMixin = {
    data() {
        return {}
    },
    methods: {
        //2006-01-02T15:04:05Z07:00  (2023-07-04T23:57:50)
        getBackendFormat(date) {
            if (this.isDateEmpty(date)) {
                return null;
            }
            return date + "T12:05:05.576147+02:00";
        },
        getBackendFormatWithTime(date, time) {
            if (this.isDateEmpty(date)) {
                return null;
            }
            return date + "T" + time + ":05.576147+02:00"
        },
        getDateInMMDDYYYYFormat(date) {
            if (this.isDateEmpty(date)) {
                return null;
            }
            return date.split('T')[0];
        },
        isToday(date) {
            const today = new Date();

            if (today.toDateString() === date.toDateString()) {
                return true;
            }

            return false;
        },
        getTime(date) {
            if (this.isDateEmpty(date)) {
                return null;
            }

            var t = date.split('T')[1];
            return t.split(':')[0] + ":" + t.split(':')[1];
        }
        ,
// formatDate(inputDate) {
//     if (typeof(inputDate) === 'string')
//     {
//         let date = inputDate.split('T')[0];
//         let day, month, year;
//
//         day = date.split('-')[2];
//         month = date.split('-')[1];
//         year = date.split('-')[0];
//
//         return `${day}-${month}-${year}`;
//     }
//
//     return inputDate.getDate() + "-" + inputDate.getMonth() + 1 + "" + inputDate.getFullYear();
// },
// formatDateWithPoints(inputDate) {
//     if (typeof(inputDate) === 'string') {
//         let date = inputDate.split('T')[0];
//         let day, month, year;
//
//         day = date.split('-')[2];
//         month = date.split('-')[1];
//         year = date.split('-')[0];
//
//         return `${day}.${month}.${year}.`;
//     }
//
//     return inputDate.getDate() + "." + inputDate.getMonth() + 1 + "." + inputDate.getFullYear() + ".";
// },
// formatDateWithPointsExceptLast(inputDate) {
//     if (typeof(inputDate) === 'string') {
//         let date = inputDate.split('T')[0];
//         let day, month, year;
//
//         day = date.split('-')[2];
//         month = date.split('-')[1];
//         year = date.split('-')[0];
//
//         return `${day}.${month}.${year}`;
//     }
//
//     var month = '' + (inputDate.getMonth() + 1);
//     var day = '' + inputDate.getDate();
//     var year = inputDate.getFullYear();
//
//     if (month.length < 2)
//         month = '0' + month;
//     if (day.length < 2)
//         day = '0' + day;
//
//     return [day, month, year].join('.');
// },
        isDateEmpty(inputDate) {
            if (inputDate === "") {
                return true;
            }
            if (inputDate === null) {
                return true;
            }
            if (typeof (inputDate) === 'object') {
                return false;
            }
            let date = inputDate.split('T')[0];
            const year = date.split('-')[0];

            return year === '0001';
        }
    },
}