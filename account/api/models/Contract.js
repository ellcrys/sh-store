"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const helpers_1 = require("./helpers");
const Promise = require("bluebird");
/**
 * Get contract
 *
 * @export
 * @param {string} email
 * @returns {Promise<any>}
 */
function get(q) {
    return new Promise((resolve, reject) => {
        return helpers_1.castDB(sails.db).models.contract.findOne({ where: q }).then(resolve).catch(reject);
    });
}
exports.get = get;
/**
 * Create a contract
 *
 * @export
 * @param {Contract} obj
 * @returns {Promise<any>}
 */
function create(obj) {
    return new Promise((resolve, reject) => {
        return helpers_1.castDB(sails.db).models.contract.create(obj).then(resolve).catch(reject);
    });
}
exports.create = create;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9tb2RlbHMvQ29udHJhY3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSx1Q0FBZ0M7QUFDaEMsTUFBTSxPQUFPLEdBQUcsT0FBTyxDQUFDLFVBQVUsQ0FBQyxDQUFBO0FBR25DOzs7Ozs7R0FNRztBQUNILGFBQW9CLENBQU07SUFDdEIsTUFBTSxDQUFDLElBQUksT0FBTyxDQUFDLENBQUMsT0FBTyxFQUFFLE1BQU07UUFDL0IsTUFBTSxDQUFDLGdCQUFNLENBQUMsS0FBSyxDQUFDLEVBQUUsQ0FBQyxDQUFDLE1BQU0sQ0FBQyxRQUFRLENBQUMsT0FBTyxDQUFDLEVBQUMsS0FBSyxFQUFFLENBQUMsRUFBQyxDQUFDLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxDQUFDLEtBQUssQ0FBQyxNQUFNLENBQUMsQ0FBQTtJQUMzRixDQUFDLENBQUMsQ0FBQTtBQUNOLENBQUM7QUFKRCxrQkFJQztBQUVEOzs7Ozs7R0FNRztBQUNILGdCQUF1QixHQUFvQjtJQUN2QyxNQUFNLENBQUMsSUFBSSxPQUFPLENBQUMsQ0FBQyxPQUFPLEVBQUUsTUFBTTtRQUMvQixNQUFNLENBQUMsZ0JBQU0sQ0FBQyxLQUFLLENBQUMsRUFBRSxDQUFDLENBQUMsTUFBTSxDQUFDLFFBQVEsQ0FBQyxNQUFNLENBQUMsR0FBRyxDQUFDLENBQUMsSUFBSSxDQUFDLE9BQU8sQ0FBQyxDQUFDLEtBQUssQ0FBQyxNQUFNLENBQUMsQ0FBQTtJQUNuRixDQUFDLENBQUMsQ0FBQTtBQUNOLENBQUM7QUFKRCx3QkFJQyIsImZpbGUiOiJhcGkvbW9kZWxzL0NvbnRyYWN0LmpzIiwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IHtjYXN0REJ9IGZyb20gJy4vaGVscGVycydcbmNvbnN0IFByb21pc2UgPSByZXF1aXJlKFwiYmx1ZWJpcmRcIilcblxuXG4vKipcbiAqIEdldCBjb250cmFjdFxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gZW1haWwgXG4gKiBAcmV0dXJucyB7UHJvbWlzZTxhbnk+fSBcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGdldChxOiBhbnkpOiBQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIHJldHVybiBjYXN0REIoc2FpbHMuZGIpLm1vZGVscy5jb250cmFjdC5maW5kT25lKHt3aGVyZTogcX0pLnRoZW4ocmVzb2x2ZSkuY2F0Y2gocmVqZWN0KVxuICAgIH0pXG59XG5cbi8qKlxuICogQ3JlYXRlIGEgY29udHJhY3RcbiAqIFxuICogQGV4cG9ydFxuICogQHBhcmFtIHtDb250cmFjdH0gb2JqIFxuICogQHJldHVybnMge1Byb21pc2U8YW55Pn0gXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBjcmVhdGUob2JqOiBtb2RlbHMuQ29udHJhY3QpOiBQcm9taXNlPGFueT4ge1xuICAgIHJldHVybiBuZXcgUHJvbWlzZSgocmVzb2x2ZSwgcmVqZWN0KSA9PiB7XG4gICAgICAgIHJldHVybiBjYXN0REIoc2FpbHMuZGIpLm1vZGVscy5jb250cmFjdC5jcmVhdGUob2JqKS50aGVuKHJlc29sdmUpLmNhdGNoKHJlamVjdClcbiAgICB9KVxufSAiXX0=
