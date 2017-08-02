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
 * Send an account confirmation message
 *
 * @export
 * @param {string} email
 * @returns {Promise<any>}
 */
function sendAccountConfirmation(to, data) {
    return new Promise((resolve, reject) => {
        try {
            Tempman.setDir(path.join(__dirname, './email_temp')).getFile(`account_confirmation.njk`, data, (err, content) => __awaiter(this, void 0, void 0, function* () {
                if (err)
                    return reject(err);
                let data = yield client.transmissions.send({
                    options: {},
                    content: {
                        from: { name: 'Ellcrys', email: 'no-reply@' + sails.config.app.emailDomain },
                        subject: 'Welcome to Ellcrys, please confirm your account',
                        html: content
                    },
                    recipients: [
                        { address: to }
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
exports.sendAccountConfirmation = sendAccountConfirmation;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9zZXJ2aWNlcy9FbWFpbFNlcnZpY2UudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7OztBQUNBLG1DQUFtQztBQUNuQyw2QkFBNkI7QUFDN0IsdUNBQXVDO0FBQ3ZDLE1BQU0sTUFBTSxHQUFHLElBQUksU0FBUyxDQUFDLEtBQUssQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLFlBQVksQ0FBQyxDQUFBO0FBRTNEOzs7Ozs7R0FNRztBQUNILGlDQUF3QyxFQUFVLEVBQUUsSUFBUztJQUN6RCxNQUFNLENBQUMsSUFBSSxPQUFPLENBQUMsQ0FBQyxPQUFPLEVBQUUsTUFBTTtRQUMvQixJQUFJLENBQUM7WUFDRCxPQUFPLENBQUMsTUFBTSxDQUFDLElBQUksQ0FBQyxJQUFJLENBQUMsU0FBUyxFQUFFLGNBQWMsQ0FBQyxDQUFDLENBQUMsT0FBTyxDQUFDLDBCQUEwQixFQUFFLElBQUksRUFBRSxDQUFPLEdBQUcsRUFBRSxPQUFPO2dCQUM5RyxFQUFFLENBQUMsQ0FBQyxHQUFHLENBQUM7b0JBQUMsTUFBTSxDQUFDLE1BQU0sQ0FBQyxHQUFHLENBQUMsQ0FBQztnQkFDNUIsSUFBSSxJQUFJLEdBQUcsTUFBTSxNQUFNLENBQUMsYUFBYSxDQUFDLElBQUksQ0FBQztvQkFDdkMsT0FBTyxFQUFFLEVBQUU7b0JBQ1gsT0FBTyxFQUFFO3dCQUNMLElBQUksRUFBRSxFQUFFLElBQUksRUFBRSxTQUFTLEVBQUUsS0FBSyxFQUFFLFdBQVcsR0FBRyxLQUFLLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxXQUFXLEVBQUU7d0JBQzVFLE9BQU8sRUFBRSxpREFBaUQ7d0JBQzFELElBQUksRUFBRSxPQUFPO3FCQUNoQjtvQkFDRCxVQUFVLEVBQUU7d0JBQ1IsRUFBRSxPQUFPLEVBQUUsRUFBRSxFQUFFO3FCQUNsQjtpQkFDSixDQUFDLENBQUE7Z0JBQ0YsT0FBTyxDQUFDLElBQUksQ0FBQyxDQUFBO1lBQ2pCLENBQUMsQ0FBQSxDQUFDLENBQUE7UUFDTixDQUFDO1FBQUMsS0FBSyxDQUFDLENBQUMsQ0FBQyxDQUFDLENBQUMsQ0FBQztZQUNULE9BQU8sQ0FBQyxLQUFLLENBQUMsNENBQTRDLEVBQUUsQ0FBQyxDQUFDLE9BQU8sQ0FBQyxDQUFBO1lBQ3RFLE1BQU0sQ0FBQyxDQUFDLENBQUMsQ0FBQTtRQUNiLENBQUM7SUFDTCxDQUFDLENBQUMsQ0FBQTtBQUNOLENBQUM7QUF2QkQsMERBdUJDIiwiZmlsZSI6ImFwaS9zZXJ2aWNlcy9FbWFpbFNlcnZpY2UuanMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgQmx1ZWJpcmQgPSByZXF1aXJlKFwiYmx1ZWJpcmRcIilcbmltcG9ydCBUZW1wbWFuID0gcmVxdWlyZShcInRlbXBtYW5cIilcbmltcG9ydCBwYXRoID0gcmVxdWlyZShcInBhdGhcIilcbmltcG9ydCBTcGFya1Bvc3QgPSByZXF1aXJlKCdzcGFya3Bvc3QnKVxuY29uc3QgY2xpZW50ID0gbmV3IFNwYXJrUG9zdChzYWlscy5jb25maWcuYXBwLnNwYXJrUG9zdEtleSlcblxuLyoqXG4gKiBTZW5kIGFuIGFjY291bnQgY29uZmlybWF0aW9uIG1lc3NhZ2VcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtzdHJpbmd9IGVtYWlsIFxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBzZW5kQWNjb3VudENvbmZpcm1hdGlvbih0bzogc3RyaW5nLCBkYXRhOiBhbnkpOiBQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIHRyeSB7XG4gICAgICAgICAgICBUZW1wbWFuLnNldERpcihwYXRoLmpvaW4oX19kaXJuYW1lLCAnLi9lbWFpbF90ZW1wJykpLmdldEZpbGUoYGFjY291bnRfY29uZmlybWF0aW9uLm5qa2AsIGRhdGEsIGFzeW5jIChlcnIsIGNvbnRlbnQpID0+IHtcbiAgICAgICAgICAgICAgICBpZiAoZXJyKSByZXR1cm4gcmVqZWN0KGVycik7XG4gICAgICAgICAgICAgICAgbGV0IGRhdGEgPSBhd2FpdCBjbGllbnQudHJhbnNtaXNzaW9ucy5zZW5kKHtcbiAgICAgICAgICAgICAgICAgICAgb3B0aW9uczoge30sXG4gICAgICAgICAgICAgICAgICAgIGNvbnRlbnQ6IHtcbiAgICAgICAgICAgICAgICAgICAgICAgIGZyb206IHsgbmFtZTogJ0VsbGNyeXMnLCBlbWFpbDogJ25vLXJlcGx5QCcgKyBzYWlscy5jb25maWcuYXBwLmVtYWlsRG9tYWluIH0sXG4gICAgICAgICAgICAgICAgICAgICAgICBzdWJqZWN0OiAnV2VsY29tZSB0byBFbGxjcnlzLCBwbGVhc2UgY29uZmlybSB5b3VyIGFjY291bnQnLFxuICAgICAgICAgICAgICAgICAgICAgICAgaHRtbDogY29udGVudFxuICAgICAgICAgICAgICAgICAgICB9LFxuICAgICAgICAgICAgICAgICAgICByZWNpcGllbnRzOiBbXG4gICAgICAgICAgICAgICAgICAgICAgICB7IGFkZHJlc3M6IHRvIH1cbiAgICAgICAgICAgICAgICAgICAgXVxuICAgICAgICAgICAgICAgIH0pXG4gICAgICAgICAgICAgICAgcmVzb2x2ZShkYXRhKVxuICAgICAgICAgICAgfSlcbiAgICAgICAgfSBjYXRjaCAoZSkge1xuICAgICAgICAgICAgY29uc29sZS5lcnJvcihcIkVycm9yIHNlbmRpbmcgYWNjb3VudCBjb25maXJtYXRpb24gbWVzc2FnZVwiLCBlLm1lc3NhZ2UpXG4gICAgICAgICAgICByZWplY3QoZSlcbiAgICAgICAgfVxuICAgIH0pXG59Il19
