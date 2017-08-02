"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const helpers_1 = require("./helpers");
const Promise = require("bluebird");
/**
 * Get account by email
 *
 * @export
 * @param {string} email
 * @returns {Promise<any>}
 */
function get(q) {
    return new Promise((resolve, reject) => {
        return helpers_1.castDB(sails.db).models.account.findOne({ where: q }).then(resolve).catch(reject);
    });
}
exports.get = get;
/**
 * Create an account
 *
 * @export
 * @param {Account} obj
 * @returns {Promise<any>}
 */
function create(obj) {
    return new Promise((resolve, reject) => {
        return helpers_1.castDB(sails.db).models.account.create(obj).then(resolve).catch(reject);
    });
}
exports.create = create;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9tb2RlbHMvQWNjb3VudC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBLHVDQUFnQztBQUNoQyxNQUFNLE9BQU8sR0FBRyxPQUFPLENBQUMsVUFBVSxDQUFDLENBQUE7QUFHbkM7Ozs7OztHQU1HO0FBQ0gsYUFBb0IsQ0FBTTtJQUN0QixNQUFNLENBQUMsSUFBSSxPQUFPLENBQUMsQ0FBQyxPQUFPLEVBQUUsTUFBTTtRQUMvQixNQUFNLENBQUMsZ0JBQU0sQ0FBQyxLQUFLLENBQUMsRUFBRSxDQUFDLENBQUMsTUFBTSxDQUFDLE9BQU8sQ0FBQyxPQUFPLENBQUMsRUFBQyxLQUFLLEVBQUUsQ0FBQyxFQUFDLENBQUMsQ0FBQyxJQUFJLENBQUMsT0FBTyxDQUFDLENBQUMsS0FBSyxDQUFDLE1BQU0sQ0FBQyxDQUFBO0lBQzFGLENBQUMsQ0FBQyxDQUFBO0FBQ04sQ0FBQztBQUpELGtCQUlDO0FBRUQ7Ozs7OztHQU1HO0FBQ0gsZ0JBQXVCLEdBQVk7SUFDL0IsTUFBTSxDQUFDLElBQUksT0FBTyxDQUFDLENBQUMsT0FBTyxFQUFFLE1BQU07UUFDL0IsTUFBTSxDQUFDLGdCQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxPQUFPLENBQUMsTUFBTSxDQUFDLEdBQUcsQ0FBQyxDQUFDLElBQUksQ0FBQyxPQUFPLENBQUMsQ0FBQyxLQUFLLENBQUMsTUFBTSxDQUFDLENBQUE7SUFDbEYsQ0FBQyxDQUFDLENBQUE7QUFDTixDQUFDO0FBSkQsd0JBSUMiLCJmaWxlIjoiYXBpL21vZGVscy9BY2NvdW50LmpzIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHtjYXN0REJ9IGZyb20gJy4vaGVscGVycydcbmNvbnN0IFByb21pc2UgPSByZXF1aXJlKFwiYmx1ZWJpcmRcIilcblxuXG4vKipcbiAqIEdldCBhY2NvdW50IGJ5IGVtYWlsXG4gKiBcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7c3RyaW5nfSBlbWFpbCBcbiAqIEByZXR1cm5zIHtQcm9taXNlPGFueT59IFxuICovXG5leHBvcnQgZnVuY3Rpb24gZ2V0KHE6IGFueSk6IFByb21pc2U8YW55PiB7XG4gICAgcmV0dXJuIG5ldyBQcm9taXNlKChyZXNvbHZlLCByZWplY3QpID0+IHtcbiAgICAgICAgcmV0dXJuIGNhc3REQihzYWlscy5kYikubW9kZWxzLmFjY291bnQuZmluZE9uZSh7d2hlcmU6IHF9KS50aGVuKHJlc29sdmUpLmNhdGNoKHJlamVjdClcbiAgICB9KVxufVxuXG4vKipcbiAqIENyZWF0ZSBhbiBhY2NvdW50XG4gKiBcbiAqIEBleHBvcnRcbiAqIEBwYXJhbSB7QWNjb3VudH0gb2JqIFxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBjcmVhdGUob2JqOiBBY2NvdW50KTogUHJvbWlzZTxhbnk+IHtcbiAgICByZXR1cm4gbmV3IFByb21pc2UoKHJlc29sdmUsIHJlamVjdCkgPT4ge1xuICAgICAgICByZXR1cm4gY2FzdERCKHNhaWxzLmRiKS5tb2RlbHMuYWNjb3VudC5jcmVhdGUob2JqKS50aGVuKHJlc29sdmUpLmNhdGNoKHJlamVjdClcbiAgICB9KVxufSAiXX0=
