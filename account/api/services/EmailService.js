"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const Tempman = require("tempman");
const path = require("path");
const SparkPost = require("sparkpost");
const client = new SparkPost(sails.config.app.sparkPostKey);
/**
 * Send a message
 *
 * @export
 * @param {{ tempFile: string, fromName: string, subject: string, to: string, data: any }} opt
 * @returns {Promise<any>}
 */
function sendMsg(opt) {
    return new Promise((resolve, reject) => {
        try {
            Tempman.setDir(path.join(__dirname, './email_temp')).getFile(`${opt.tempFile}.njk`, opt.data, (err, content) => __awaiter(this, void 0, void 0, function* () {
                if (err)
                    return reject(err);
                let data = yield client.transmissions.send({
                    options: {},
                    content: {
                        from: { name: opt.fromName, email: 'no-reply@' + sails.config.app.emailDomain },
                        subject: opt.subject,
                        html: content
                    },
                    recipients: [
                        { address: opt.to }
                    ]
                });
                resolve(data);
            }));
        }
        catch (e) {
            console.error("Error sending account confirmation message", e.message);
            reject(e);
        }
    });
}
exports.sendMsg = sendMsg;
/**
 * Send an account confirmation message
 *
 * @export
 * @param {string} email
 * @returns {Promise<any>}
 */
function sendAccountConfirmation(to, data) {
    return sendMsg({
        tempFile: 'account_confirmation',
        fromName: 'Ellcrys',
        subject: 'Welcome to Ellcrys, please confirm your account',
        to,
        data,
    });
}
exports.sendAccountConfirmation = sendAccountConfirmation;
/**
 * Send a password reset email
 *
 * @export
 * @param {string} email
 * @returns {Promise<any>}
 */
function sendPasswordResetEmail(to, data) {
    return sendMsg({
        tempFile: 'reset-password',
        fromName: 'Ellcrys',
        subject: 'Reset your Ellcrys password',
        to,
        data,
    });
}
exports.sendPasswordResetEmail = sendPasswordResetEmail;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9zZXJ2aWNlcy9FbWFpbFNlcnZpY2UudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7OztBQUNBLG1DQUFtQztBQUNuQyw2QkFBNkI7QUFDN0IsdUNBQXVDO0FBQ3ZDLE1BQU0sTUFBTSxHQUFHLElBQUksU0FBUyxDQUFDLEtBQUssQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLFlBQVksQ0FBQyxDQUFBO0FBRTNEOzs7Ozs7R0FNRztBQUNILGlCQUF3QixHQUFtRjtJQUN2RyxNQUFNLENBQUMsSUFBSSxPQUFPLENBQUMsQ0FBQyxPQUFPLEVBQUUsTUFBTTtRQUMvQixJQUFJLENBQUM7WUFDRCxPQUFPLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsU0FBUyxFQUFFLGNBQWMsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDLEdBQUcsR0FBRyxDQUFDLFFBQVEsTUFBTSxFQUFFLEdBQUcsQ0FBQyxJQUFJLEVBQUUsQ0FBTyxHQUFHLEVBQUUsT0FBTztnQkFDN0csRUFBRSxDQUFDLENBQUMsR0FBRyxDQUFDO29CQUFDLE1BQU0sQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLENBQUM7Z0JBQzVCLElBQUksSUFBSSxHQUFHLE1BQU0sTUFBTSxDQUFDLGFBQWEsQ0FBQyxJQUFJLENBQUM7b0JBQ3ZDLE9BQU8sRUFBRSxFQUFFO29CQUNYLE9BQU8sRUFBRTt3QkFDTCxJQUFJLEVBQUUsRUFBRSxJQUFJLEVBQUUsR0FBRyxDQUFDLFFBQVEsRUFBRSxLQUFLLEVBQUUsV0FBVyxHQUFHLEtBQUssQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLFdBQVcsRUFBRTt3QkFDL0UsT0FBTyxFQUFFLEdBQUcsQ0FBQyxPQUFPO3dCQUNwQixJQUFJLEVBQUUsT0FBTztxQkFDaEI7b0JBQ0QsVUFBVSxFQUFFO3dCQUNSLEVBQUUsT0FBTyxFQUFFLEdBQUcsQ0FBQyxFQUFFLEVBQUU7cUJBQ3RCO2lCQUNKLENBQUMsQ0FBQTtnQkFDRixPQUFPLENBQUMsSUFBSSxDQUFDLENBQUE7WUFDakIsQ0FBQyxDQUFBLENBQUMsQ0FBQTtRQUNOLENBQUM7UUFBQyxLQUFLLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQ1QsT0FBTyxDQUFDLEtBQUssQ0FBQyw0Q0FBNEMsRUFBRSxDQUFDLENBQUMsT0FBTyxDQUFDLENBQUE7WUFDdEUsTUFBTSxDQUFDLENBQUMsQ0FBQyxDQUFBO1FBQ2IsQ0FBQztJQUNMLENBQUMsQ0FBQyxDQUFBO0FBQ04sQ0FBQztBQXZCRCwwQkF1QkM7QUFFRDs7Ozs7O0dBTUc7QUFDSCxpQ0FBd0MsRUFBVSxFQUFFLElBQVM7SUFDekQsTUFBTSxDQUFDLE9BQU8sQ0FBQztRQUNYLFFBQVEsRUFBRSxzQkFBc0I7UUFDaEMsUUFBUSxFQUFFLFNBQVM7UUFDbkIsT0FBTyxFQUFFLGlEQUFpRDtRQUMxRCxFQUFFO1FBQ0YsSUFBSTtLQUNQLENBQUMsQ0FBQTtBQUNOLENBQUM7QUFSRCwwREFRQztBQUVEOzs7Ozs7R0FNRztBQUNILGdDQUF1QyxFQUFVLEVBQUUsSUFBUztJQUN4RCxNQUFNLENBQUMsT0FBTyxDQUFDO1FBQ1gsUUFBUSxFQUFFLGdCQUFnQjtRQUMxQixRQUFRLEVBQUUsU0FBUztRQUNuQixPQUFPLEVBQUUsNkJBQTZCO1FBQ3RDLEVBQUU7UUFDRixJQUFJO0tBQ1AsQ0FBQyxDQUFBO0FBQ04sQ0FBQztBQVJELHdEQVFDIiwiZmlsZSI6ImFwaS9zZXJ2aWNlcy9FbWFpbFNlcnZpY2UuanMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgQmx1ZWJpcmQgPSByZXF1aXJlKFwiYmx1ZWJpcmRcIilcbmltcG9ydCBUZW1wbWFuID0gcmVxdWlyZShcInRlbXBtYW5cIilcbmltcG9ydCBwYXRoID0gcmVxdWlyZShcInBhdGhcIilcbmltcG9ydCBTcGFya1Bvc3QgPSByZXF1aXJlKCdzcGFya3Bvc3QnKVxuY29uc3QgY2xpZW50ID0gbmV3IFNwYXJrUG9zdChzYWlscy5jb25maWcuYXBwLnNwYXJrUG9zdEtleSlcblxuLyoqXG4gKiBTZW5kIGEgbWVzc2FnZVxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3sgdGVtcEZpbGU6IHN0cmluZywgZnJvbU5hbWU6IHN0cmluZywgc3ViamVjdDogc3RyaW5nLCB0bzogc3RyaW5nLCBkYXRhOiBhbnkgfX0gb3B0IFxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBzZW5kTXNnKG9wdDogeyB0ZW1wRmlsZTogc3RyaW5nLCBmcm9tTmFtZTogc3RyaW5nLCBzdWJqZWN0OiBzdHJpbmcsIHRvOiBzdHJpbmcsIGRhdGE6IGFueSB9KTogUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICB0cnkge1xuICAgICAgICAgICAgVGVtcG1hbi5zZXREaXIocGF0aC5qb2luKF9fZGlybmFtZSwgJy4vZW1haWxfdGVtcCcpKS5nZXRGaWxlKGAke29wdC50ZW1wRmlsZX0ubmprYCwgb3B0LmRhdGEsIGFzeW5jIChlcnIsIGNvbnRlbnQpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoZXJyKSByZXR1cm4gcmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgbGV0IGRhdGEgPSBhd2FpdCBjbGllbnQudHJhbnNtaXNzaW9ucy5zZW5kKHtcbiAgICAgICAgICAgICAgICAgICAgb3B0aW9uczoge30sXG4gICAgICAgICAgICAgICAgICAgIGNvbnRlbnQ6IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGZyb206IHsgbmFtZTogb3B0LmZyb21OYW1lLCBlbWFpbDogJ25vLXJlcGx5QCcgKyBzYWlscy5jb25maWcuYXBwLmVtYWlsRG9tYWluIH0sXG4gICAgICAgICAgICAgICAgICAgICAgICBzdWJqZWN0OiBvcHQuc3ViamVjdCxcbiAgICAgICAgICAgICAgICAgICAgICAgIGh0bWw6IGNvbnRlbnRcbiAgICAgICAgICAgICAgICAgICAgfSxcbiAgICAgICAgICAgICAgICAgICAgcmVjaXBpZW50czogW1xuICAgICAgICAgICAgICAgICAgICAgICAgeyBhZGRyZXNzOiBvcHQudG8gfVxuICAgICAgICAgICAgICAgICAgICBdXG4gICAgICAgICAgICAgICAgfSlcbiAgICAgICAgICAgICAgICByZXNvbHZlKGRhdGEpXG4gICAgICAgICAgICB9KVxuICAgICAgICB9IGNhdGNoIChlKSB7XG4gICAgICAgICAgICBjb25zb2xlLmVycm9yKFwiRXJyb3Igc2VuZGluZyBhY2NvdW50IGNvbmZpcm1hdGlvbiBtZXNzYWdlXCIsIGUubWVzc2FnZSlcbiAgICAgICAgICAgIHJlamVjdChlKVxuICAgICAgICB9XG4gICAgfSlcbn1cblxuLyoqXG4gKiBTZW5kIGFuIGFjY291bnQgY29uZmlybWF0aW9uIG1lc3NhZ2VcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGVtYWlsIFxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBzZW5kQWNjb3VudENvbmZpcm1hdGlvbih0bzogc3RyaW5nLCBkYXRhOiBhbnkpOiBQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBzZW5kTXNnKHtcbiAgICAgICAgdGVtcEZpbGU6ICdhY2NvdW50X2NvbmZpcm1hdGlvbicsXG4gICAgICAgIGZyb21OYW1lOiAnRWxsY3J5cycsXG4gICAgICAgIHN1YmplY3Q6ICdXZWxjb21lIHRvIEVsbGNyeXMsIHBsZWFzZSBjb25maXJtIHlvdXIgYWNjb3VudCcsXG4gICAgICAgIHRvLFxuICAgICAgICBkYXRhLFxuICAgIH0pXG59XG5cbi8qKlxuICogU2VuZCBhIHBhc3N3b3JkIHJlc2V0IGVtYWlsXG4gKiBcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBlbWFpbCBcbiAqIEByZXR1cm5zIHtQcm9taXNlPGFueT59IFxuICovXG5leHBvcnQgZnVuY3Rpb24gc2VuZFBhc3N3b3JkUmVzZXRFbWFpbCh0bzogc3RyaW5nLCBkYXRhOiBhbnkpOiBQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBzZW5kTXNnKHtcbiAgICAgICAgdGVtcEZpbGU6ICdyZXNldC1wYXNzd29yZCcsXG4gICAgICAgIGZyb21OYW1lOiAnRWxsY3J5cycsXG4gICAgICAgIHN1YmplY3Q6ICdSZXNldCB5b3VyIEVsbGNyeXMgcGFzc3dvcmQnLFxuICAgICAgICB0byxcbiAgICAgICAgZGF0YSxcbiAgICB9KVxufSJdfQ==
