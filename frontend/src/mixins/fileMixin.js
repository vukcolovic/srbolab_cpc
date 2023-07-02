// import axios from "axios";
// import notie from 'notie';

export const fileMixin = {
    data() {
        return {

        }
    },
    methods: {
        base64ToArrayBuffer(base64) {
            var binaryString = window.atob(base64);
            var binaryLen = binaryString.length;
            var bytes = new Uint8Array(binaryLen);
            for (var i = 0; i < binaryLen; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }
            return bytes;
        }
    },
}